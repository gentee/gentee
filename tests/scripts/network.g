const {
   README = `https://github.com/gentee/gentee`
   TITLE = `Gentee programming language`
   TESTURL = `https://..../`
}

func testrequest() {
    map empty head
    str response = HTTPRequest(TESTURL + "get", "GET", empty,head)
    Println(response)
    map params = {  `par1`: `Имя`, `value_2`: `This is a string`  }
    Println(HTTPRequest(TESTURL + "get", "GET", params, head))
    Println(HTTPRequest(TESTURL + "post", "POST", params, head))
    map headjson = { `Content-Type`: `application/json; charset=UTF-8` }
    Println(HTTPRequest(TESTURL + "post", "POST", params, headjson))
}

run str {
//    testrequest()
    hinfo hi = HeadInfo(`https://github.com/gentee/gentee/`)
    if hi.Status!= 200 : error(103,`HeadInfo status`)
    str page = HTTPPage(README)
    if Find(page, TITLE) < 0 : error(100, `HTTPPage`)
    buf bufPage = HTTPGet(README)
    if Find(str(bufPage), TITLE) < 0 : error(101, `HTTPGet`)
    str ftemp = TempDir() + `/readme.html`
    int size = Download(README, ftemp)
    if Max(*page, size) - Min(*page, size) > size/100 : error(102, `Download`)
    Remove(ftemp)
    return `OK`
}