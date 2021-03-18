package pagination

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//Pagination 分页器
type Pagination struct {
	Request *http.Request
	Total   int //总行数
	Pagesize int //每页的数量
}

//Initialize 初始化分页器
func Initialize(req *http.Request, total int, pagesize int) *Pagination {
	return &Pagination{
		Request: req,
		Total:   total,
		Pagesize: pagesize,
	}
}

//Pages 渲染生成html分页标签
func (p *Pagination) Pages() string {
	queryParams := p.Request.URL.Query()
	//从当前请求中获取page
	page := queryParams.Get("page")
	if page == "" {page = "1"}
	//将页码转换成整型，以便计算
	currentPage, _ := strconv.Atoi(page)
	if currentPage == 0 {return ""}

	//计算总页数
	var totalPageNum = int(math.Ceil(float64(p.Total) / float64(p.Pagesize)))

	//首页链接
	var firstLink string
	//上一页链接
	var prevLink string
	//下一页链接
	var nextLink string
	//末页链接
	var lastLink string
	//中间页码链接
	var pageLinks []string

	//当总页数小于等于1时，不返回分页的html
	if totalPageNum <= 1 {return ""}

	//首页和上一页链接
	if currentPage > 1 {
		firstLink = fmt.Sprintf(`<li class="page-item">
		<a class="page-link" href="%s" aria-label="Previous">
		<span aria-hidden="true">首页</span>
		</a>
		</li>`, p.pageURL("1"))

		prevLink = fmt.Sprintf(`<li class="page-item">
		<a class="page-link" href="%s" aria-label="Previous">
		<span aria-hidden="true">上一页</span>
		</a>
		</li>`, p.pageURL(strconv.Itoa(currentPage-1)))
	}

	//末页和下一页
	if currentPage < totalPageNum {
		lastLink = fmt.Sprintf(`<li class="page-item">
                <a class="page-link" href="%s" aria-label="Next">
                    <span aria-hidden="true">末页</span>
                </a>
            </li>`, p.pageURL(strconv.Itoa(totalPageNum)))
		nextLink = fmt.Sprintf(`<li class="page-item">
                <a class="page-link" href="%s" aria-label="Next">
                    <span aria-hidden="true">下一页</span>
                </a>
            </li>`, p.pageURL(strconv.Itoa(currentPage+1)))
	}

	//生成中间页码链接
	pageLinks = make([]string, 0, 10)
	startPos := currentPage - 3
	endPos := currentPage + 3
	if startPos < 1 {
		endPos = endPos + int(math.Abs(float64(startPos))) + 1
		startPos = 1
	}
	if endPos > totalPageNum {
		endPos = totalPageNum
	}
	for i := startPos; i <= endPos; i++ {
		var s string
		if i == currentPage {
			s = fmt.Sprintf(`<li class="page-item active"><a class="page-link" href="%s">%d</a></li>`, p.pageURL(strconv.Itoa(i)), i)
		} else {
			s = fmt.Sprintf(`<li class="page-item"><a class="page-link" href="%s">%d</a></li>`, p.pageURL(strconv.Itoa(i)), i)
		}
		pageLinks = append(pageLinks, s)
	}

	return fmt.Sprintf(`<nav aria-label="Page navigation ">
	<ul class="pagination">%s%s%s%s%s</ul>
	</nav>`, firstLink, prevLink, strings.Join(pageLinks, ""), nextLink, lastLink)
}

//pageURL 生成分页url
func (p *Pagination) pageURL(page string) string {
	//基于当前url新建一个url对象
	u, _ := url.Parse(p.Request.URL.String()) // /articles?page=2&id=10
	q := u.Query()                            //map[page:[2] id:[10]]
	q.Set("page", page)
	u.RawQuery = q.Encode()
	return u.String()
}
