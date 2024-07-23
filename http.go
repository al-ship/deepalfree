package deepalfree

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type httpServ struct {
	srv             *http.Server
	filesServerHost string
	apps            []App
	logger          *zap.Logger
	path            string
}

func NewHttp(filesServHost string, logger *zap.Logger) *httpServ {
	r := mux.NewRouter()
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	s := &httpServ{
		srv: &http.Server{
			Handler: r,
			Addr:    "0.0.0.0:7070",
		},
		filesServerHost: filesServHost,
		logger:          logger,
		path:            dir,
	}
	appStore := r.PathPrefix("/hu-apigw/appstore").Subrouter()
	appStore.HandleFunc("/api/v1/app/list", s.appList).Methods("GET")
	appStore.HandleFunc("/api/v1/task/initial-params", s.initialParams).Methods("GET")
	r.HandleFunc("/hu-apigw/wiki/api/v1/commodity/purchase-list", s.purchaseList).Methods("GET")
	r.HandleFunc("/hu-apigw/wiki/api/v1/commodity/resource-list", s.resourceList).Methods("GET")

	r.PathPrefix("/hu-apigw/apps/").Handler(http.StripPrefix("/hu-apigw/apps/", http.FileServer(http.Dir(s.path+"/apps"))))
	return s
}

func (h *httpServ) Start() error {
	files, err := os.ReadDir(h.path + "/apps")
	if err != nil {
		return fmt.Errorf("read apps files failed: %w", err)
	}
	appID := 0
	h.logger.Info("searching apps...")
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		appID++
		fPath := h.path + "/apps/" + file.Name()
		content, err := os.ReadFile(fPath)
		if err != nil {
			return fmt.Errorf("read file %s failed: %w", file.Name(), err)
		}
		sum := md5.Sum(content)

		h.apps = append(h.apps, App{
			ApkInfo: ApkInfo{
				ApkName:  file.Name(),
				FileSize: fmt.Sprintf("%d", len(content)),
				HashCode: hex.EncodeToString(sum[:]),
				HashType: "md5",
				URL:      fmt.Sprintf("https://%s/hu-apigw/apps/%s", h.filesServerHost, file.Name()),
			},
			AppID:       strconv.Itoa(appID),
			Icon:        "",
			Name:        file.Name(),
			PackageName: "app",
			PayInfo: PayInfo{
				Discount:      100.0,
				OrderSource:   "APPSTORE",
				OriginalPrice: 0.0,
				PayWays:       []any{},
				Price:         0.0,
			},
			RestrictedState: 0,
			Slogan:          file.Name(),
			Statics: Statics{
				Downloads:  "0",
				Installs:   "0",
				Uninstalls: "0",
				Updates:    "0",
			},
			Tags: []Tags{{
				TagCode:  "yes",
				TypeCode: "tuijian",
			}},
			Tid:           strconv.Itoa(appID),
			Type:          "7",
			Uninstall:     true,
			Version:       "1.0.0",
			VersionID:     strconv.Itoa(appID),
			VersionNumber: 1,
		})
		h.logger.Info("app found", zap.String("app", file.Name()))
	}
	h.logger.Sugar().Infof("found %d apps", len(h.apps))
	return h.srv.ListenAndServe()
}

func (h *httpServ) Stop() {
	if err := h.srv.Close(); err != nil {
		h.logger.Error("http server close failed", zap.Error(err))
	}
}

func (h *httpServ) appList(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", " application/json;charset=UTF-8")
	content, err := json.Marshal(&Response{
		Code: 0,
		Data: Data{
			List: h.apps,
			Pagination: Pagination{
				Current:  1,
				PageSize: 50,
				Total:    len(h.apps),
			},
		},
		Msg:     "",
		Success: true,
	})
	if err != nil {
		h.logger.Error("error marshal response", zap.Error(err))
		return
	}
	_, err = w.Write(content)
	if err != nil {
		h.logger.Error("error write applist content", zap.Error(err))
	}
}

func (h *httpServ) initialParams(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", " application/json;charset=UTF-8")
	res := `{"code":0,"data":{"effective_time":"2592000","max_tasks":10,"page_size":10,"simultaneously_downloads":1,"task_interval":5.0},"msg":"","success":true}`
	_, err := w.Write([]byte(res))
	if err != nil {
		h.logger.Error("error write initialParams content", zap.Error(err))
	}
}

func (h *httpServ) purchaseList(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", " application/json;charset=UTF-8")
	res := `{"code":0,"data":[{"purchase":false,"tid":"1689833216308191232"},{"purchase":false,"tid":"1753235507531247616"},{"purchase":false,"tid":"1731638050174312448"},{"purchase":false,"tid":"1731641971139805184"},{"purchase":false,"tid":"1730488044809453568"},{"purchase":false,"tid":"1752589759147397120"},{"purchase":false,"tid":"1689822035781619712"},{"purchase":false,"tid":"1689837073524682752"},{"purchase":false,"tid":"1689840784481247232"},{"purchase":false,"tid":"1731557639117991936"},{"purchase":false,"tid":"1731513821424537600"},{"purchase":false,"tid":"1730406973070020608"},{"purchase":false,"tid":"1731646006640422912"},{"purchase":false,"tid":"1730510939409879040"},{"purchase":false,"tid":"1689827664007794688"},{"purchase":false,"tid":"1689830619482185728"},{"purchase":false,"tid":"1689159805103017984"},{"purchase":false,"tid":"1689838936303792128"},{"purchase":false,"tid":"1730515022307262464"},{"purchase":false,"tid":"1731555708693401600"},{"purchase":false,"tid":"1752624925300891648"},{"purchase":false,"tid":"1731560057677672448"},{"purchase":false,"tid":"1752971215286870016"},{"purchase":false,"tid":"1689832388338642944"},{"purchase":false,"tid":"1731633856914505728"}],"msg":"","success":true}`
	_, err := w.Write([]byte(res))
	if err != nil {
		h.logger.Error("error write initialParams content", zap.Error(err))
	}
}

func (h *httpServ) resourceList(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", " application/json;charset=UTF-8")
	res := `{"code":0,"data":{"list":[{"action":{"approach":"jump_applyDetail","style":"INNER_JUMP","tid":"1689837075303608320"},"block":{"bid":"1763026107200602113","name":"YYBANNER!C673-EVE"},"forum":{"fid":"1698955790128029696","fup":"1698949986243715072","name":"应用商店banner"},"post":{"files":[{"dsshandle":"3352902f-be96-42fb-a228-a41e260e0e73","ext":"jpg","fileId":"1763046846704951296","md5":"08a6f52722e918a07c3f2e3be647d858","name":"优酷视频.jpg","sha256":"1b99d110b93a47e0d01b6f4eeaf42d7b70ab3d219c4c07c370a88eeaf01118b5","size":"182486","url":"https://incallcdn.changan.com.cn/static/cdnshamdown/prod/501744c9798821c506c3448832e77a41e900a4991f5ad9d2ad05eadf9ac71b38c2554082171249dca4bc49abea4827da"}],"message":"","pid":"1763047108012949504","tags":"[]"},"startdate":"2024-02-29 00:00:00","thread":{"attribute":"{}","dateline":1709178161000,"fid":"1698955790128029696","lastposter":"202313727","pics":["https://incallcdn.changan.com.cn/static/cdnshamdown/prod/501744c9798821c506c3448832e77a41e900a4991f5ad9d2ad05eadf9ac71b38c2554082171249dca4bc49abea4827da"],"recommends":0.0,"relate_by_tag":"[]","subject":"优酷视频","type_id":"0"},"threadExt":{"rangeOfChoice":1},"tid":"1763047108012949504"},{"action":{"approach":"jump_applyDetail","style":"INNER_JUMP","tid":"1689822038634004480"},"block":{"bid":"1763026107200602113","name":"YYBANNER!C673-EVE"},"forum":{"fid":"1698955790128029696","fup":"1698949986243715072","name":"应用商店banner"},"post":{"files":[{"dsshandle":"c1097945-7909-4f47-b1b6-417bce653125","ext":"jpg","fileId":"1763010836578738176","md5":"d970204d45ef143a6084017f5464afb2","name":"酷我音乐.jpg","sha256":"2408be16b5574ad48ea1ed08434cb5b5b3c418b3b554541b322341fb7b8ef6b0","size":"112938","url":"https://incallcdn.changan.com.cn/static/cdnshamdown/prod/ef9f3d840de129c94522206d1e4792dd2606ab080042aea002c9ac040d1b3340f9b059f638af708eb2c5f1423de0c261"}],"message":"","pid":"1763026007249027072","tags":"[]"},"startdate":"2024-02-29 09:00:00","thread":{"attribute":"{}","dateline":1709173130000,"fid":"1698955790128029696","lastposter":"202313727","pics":["https://incallcdn.changan.com.cn/static/cdnshamdown/prod/ef9f3d840de129c94522206d1e4792dd2606ab080042aea002c9ac040d1b3340f9b059f638af708eb2c5f1423de0c261"],"recommends":0.0,"relate_by_tag":"[]","subject":"酷我音乐385","type_id":"0"},"threadExt":{"rangeOfChoice":1},"tid":"1763026007249027072"}],"pagination":{"current":1,"pageSize":5,"total":2}},"msg":"","success":true}`
	_, err := w.Write([]byte(res))
	if err != nil {
		h.logger.Error("error write initialParams content", zap.Error(err))
	}
}
