const {
   README = `https://github.com/gentee/gentee`
   TITLE = `Gentee programming language`
}

run str {
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