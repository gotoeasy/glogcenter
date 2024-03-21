package controller

import (
	"glc/com"
	"glc/conf"
	"glc/gweb"
	"glc/ldb"
	"glc/ldb/search"
	"glc/ldb/sysmnt"
	"time"

	"github.com/gotoeasy/glang/cmn"
)

type storageItem struct {
	storeName     string // 日志仓
	total         uint32 // 日志件数
	isSearchRange bool   // 是否条件范围的日志仓
}

var cacheStoreNames []string // 所有的日志仓（避免每次读磁盘，适当使用缓存）
var cacheTime time.Time      // 最近一次读日志仓目录的时间点

// 日志检索（表单提交方式）
func LogSearchController(req *gweb.HttpRequest) *gweb.HttpResult {

	token := req.GetToken()
	username := ""
	if (!InWhiteList(req) && InBlackList(req)) || (conf.IsEnableLogin() && GetUsernameByToken(token) == "") {
		return gweb.Error403() // 黑名单检查、登录检查
	}
	if conf.IsEnableLogin() {
		username = GetUsernameByToken(token)
		catchSession.Set(token, username) // 会话重新计时
	}
	SetOrigin(req)

	// 准备好各种场景的检索条件（系统【~】、日志级别【!】、用户【@】）
	startTime := time.Now()
	mnt := sysmnt.NewSysmntStorage()
	cond := &search.SearchCondition{SearchSize: conf.GetPageSize()}
	cond.StoreName = req.GetFormParameter("storeName")                        // 日志仓条件
	cond.SearchKey = req.GetFormParameter("searchKey")                        // 输入的查询关键词
	cond.CurrentStoreName = req.GetFormParameter("currentStoreName")          // 滚动查询时定位用日志仓
	cond.CurrentId = cmn.StringToUint32(req.GetFormParameter("currentId"), 0) // 滚动查询时定位用ID
	cond.Forward = cmn.StringToBool(req.GetFormParameter("forward"), true)    // 是否向下滚动查询
	cond.DatetimeFrom = req.GetFormParameter("datetimeFrom")                  // 日期范围（From）
	cond.DatetimeTo = req.GetFormParameter("datetimeTo")                      // 日期范围（To）
	cond.OrgSystem = cmn.Trim(req.GetFormParameter("system"))                 // 系统
	cond.User = cmn.ToLower(cmn.Trim(req.GetFormParameter("user")))           // 用户
	cond.Loglevel = cmn.ToLower(req.GetFormParameter("loglevel"))             // 单选条件
	cond.Loglevels = cmn.Split(cond.Loglevel, ",")                            // 多选条件
	if cond.User != "" {
		cond.User = "@" + cond.User // 有指定用户条件
	}
	if len(cond.Loglevels) <= 1 || len(cond.Loglevels) >= 4 {
		cond.Loglevels = make([]string, 0) // 多选的单选或全选，都清空（单选走loglevel索引，全选等于没选）
	}
	if !cmn.IsBlank(cond.Loglevel) && !cmn.Contains(cond.Loglevel, ",") {
		cond.Loglevel = "!" + cmn.Trim(cond.Loglevel) // 编辑日志级别单选条件，以便精确匹配
	} else {
		cond.Loglevel = "" // 清空日志级别单选条件，以便多选配配（改用loglevels）
	}
	if !conf.IsEnableLogin() {
		cond.OrgSystems = append(cond.OrgSystems, "*") // 不需登录时全部系统都有访问权限
		if cond.OrgSystem != "" {
			cond.OrgSystem = "~" + cond.OrgSystem // 多个系统权限，按输入的系统作条件
		}
	} else {
		if username == conf.GetUsername() {
			// 管理员，不限系统
			cond.OrgSystems = append(cond.OrgSystems, "*")
			if cond.OrgSystem != "" {
				cond.OrgSystem = "~" + cond.OrgSystem // 多个系统权限，按输入的系统作条件
			}
		} else {
			// 一般用户，按设定权限
			user := mnt.GetSysUser(username)
			if user == nil {
				return gweb.Error403() // 可能出现，用户登录使用期间被管理员删除账号
			}
			if user.Systems == "*" {
				cond.OrgSystems = append(cond.OrgSystems, "*") // 全部系统都有访问权限
				if cond.OrgSystem != "" {
					cond.OrgSystem = "~" + cond.OrgSystem // 多个系统权限，按输入的系统作条件
				}
			} else {
				ary := cmn.Split(cmn.ToLower(user.Systems), ",")
				okSystem := false
				for i := 0; i < len(ary); i++ {
					cond.OrgSystems = append(cond.OrgSystems, "~"+ary[i]) // 仅设定的系统有访问权限
					if cond.OrgSystem == "" || cmn.EqualsIngoreCase(cond.OrgSystem, ary[i]) {
						okSystem = true
					}
				}

				if !okSystem {
					cond.OrgSystem = "-" // 输入的是没权限的系统时，写个不存在的条件用于快速返回
				} else {
					if len(ary) == 1 {
						cond.OrgSystem = "~" + ary[0] // 一共就一个系统权限，直接作为条件即可
					} else if cond.OrgSystem != "" {
						cond.OrgSystem = "~" + cond.OrgSystem // 多个系统权限，按输入的系统作条件
					}
				}
			}
		}
	}

	// 范围内的日志仓都查一遍
	// 注1）日志不断新增时，总件数可能会因为时间点原因不适最新，从而变现出点点小误差【完全可接受】
	// 注2）跨仓检索时，非本次检索的目标仓的话，只查取相关件数不做真正筛选计数以提高性能，最大匹配件数有时可能出现较大误差【折中可接受】
	result := &search.SearchResult{PageSize: cmn.IntToString(conf.GetPageSize())}
	var total uint32
	var count uint32
	storeItems := getStoreItems(cond.StoreName, cond.DatetimeFrom, cond.DatetimeTo)
	for i, max := 0, len(storeItems); i < max; i++ {
		item := storeItems[i]
		if !item.isSearchRange {
			// 不需要查数据，只查关联件数
			total += mnt.GetStorageDataCount(item.storeName) // 累加总件数
			continue
		}

		cond.SearchSize = conf.GetPageSize() - len(result.Data) // 本次需要查多少件
		if cond.CurrentStoreName != "" && item.storeName > cond.CurrentStoreName {
			cond.SearchSize = 0 // 是范围内的日志仓，但不是本次要查的，设为0不查数据，只查关联件数
		}

		eng := ldb.NewEngine(item.storeName)     // 遍历日志仓检索
		rs := eng.Search(cond)                   // 【检索】按动态的要求件数检索
		total += cmn.StringToUint32(rs.Total, 0) // 累加总件数
		count += cmn.StringToUint32(rs.Count, 0) // 累加最大匹配件数
		if len(rs.Data) > 0 {
			result.Data = append(result.Data, rs.Data...) // 累加查询结果
			result.LastStoreName = item.storeName         // 设定检索结果最后一条（最久远）日志所在的日志仓，页面向下滚动继续检索时定位用
		}

		if !(cond.CurrentStoreName != "" && item.storeName > cond.CurrentStoreName) {
			// 仅针对更久远的日志仓
			if len(result.Data) < conf.GetPageSize() && i < max-1 {
				// 数据没查够，且后面还有日志仓待查询，准备好跨仓查询条件
				cond.CurrentId = 0         // 下一日志仓从头开始查
				cond.CurrentStoreName = "" // 从头开始所以这个条件不再适用，清空
			}
		}
	}

	// 返回结果
	result.Total = cmn.Uint32ToString(total)                                      // 总件数
	result.Count = cmn.Uint32ToString(count)                                      // 最大匹配检索（笼统，在最大查取件数（5000件）内查完时，前端会改成精确的和结果一样的件数）
	result.TimeMessage = "耗时" + getTimeInfo(time.Since(startTime).Milliseconds()) // 查询耗时
	return gweb.Result(result)
}

// 筛选出日志仓检索范围
func getStoreItems(storeName string, datetimeFrom string, datetimeTo string) []*storageItem {
	sysmntStore := sysmnt.NewSysmntStorage()
	var items []*storageItem
	if !conf.IsStoreNameAutoAddDate() {
		// 单日志仓
		name := com.GeyStoreNameByDate("")
		items = append(items, &storageItem{storeName: name, total: sysmntStore.GetStorageDataCount(name), isSearchRange: true})
		return items
	}

	// 遍历日志仓，比较日期范围筛选日志仓
	hasDateCond := (datetimeFrom != "" && datetimeTo != "")     // 是否有日期范围条件
	from := cmn.ReplaceAll(cmn.Left(datetimeFrom, 10), "-", "") // yyyymmdd或“”
	to := cmn.ReplaceAll(cmn.Left(datetimeTo, 10), "-", "")     // yyyymmdd或“”
	if time.Since(cacheTime) >= time.Second*10 {
		cacheStoreNames = com.GetStorageNames(conf.GetStorageRoot(), ".sysmnt") // 所有的日志仓，结果已排序，缓存10秒避免频繁读盘
		cacheTime = time.Now()
	}
	for i, max := 0, len(cacheStoreNames); i < max; i++ {
		name := cacheStoreNames[i]
		item := &storageItem{storeName: name, total: sysmntStore.GetStorageDataCount(name)}
		date := cmn.Right(name, 8) // yyyymmdd
		if storeName == "" {
			// 日志仓条件空白
			if hasDateCond {
				if date >= from && date <= to {
					item.isSearchRange = true // 日期范围内的日志仓都是条件范围
				}
			} else {
				item.isSearchRange = true // 无日志仓条件、且无日期条件，全部都是条件范围了
			}
		} else {
			// 有日志仓条件
			if hasDateCond {
				if storeName == name && date >= from && date <= to {
					item.isSearchRange = true // 有日期条件，得满足日期条件，该日志仓才是条件范围
				}
			} else {
				if storeName == name {
					item.isSearchRange = true // 没日期条件，仅该日志仓是条件范围
				}
			}
		}
		items = append(items, item)
	}
	return items
}

func getTimeInfo(milliseconds int64) string {
	if milliseconds >= 1000 {
		return " " + cmn.Float64ToString(cmn.Round1(float64(milliseconds)/1000.0)) + " 秒"
	}
	return " " + cmn.Int64ToString(milliseconds) + " 毫秒"
}
