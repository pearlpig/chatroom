$(function() {
    addPage()
    show(1)
    $.ajax({
        type: "GET",
        url: "/check",
        success: function(data) {
            result = JSON.parse(data)
            console.log(result)
                // $("#memberInfo").hide()
            if (result.member_id != -1) {
                memberBar(result.nickname)
            }

        }
    });
    // signup
    $("#toSignupPage").click(function() {
        console.log("signup")
        location.href = "/signup"
    })

    // login
    $("#toLoginPage").click(function() {
        console.log("login")
        location.href = "/login"
    })

    // page button
    $('a[class*="pageBtn"]').click(function() {
        let page = parseInt(this.name)
        show(page)
    })
    $('a[class*="first"]').click(function() {
        show(1)
    })
    $('a[class="last"]').click(function() {
        show(10)
    })
    $('a[class="prev"]').click(function() {
        // var nowPage = $('a[class*="prev"]').parent().next()
        let nowPage
        for (nowPage = $('a[class*="prev"]').parent().next(); nowPage.attr('style'); nowPage = nowPage.next()) {

        }
        nowPage = parseInt(nowPage.attr('name'))
        if (nowPage == 1) {
            show(1)

        } else {
            show(nowPage - 1)
        }
    })
    $('a[class*="next"]').click(function() {
        let nowPage
        for (nowPage = $('a[class*="prev"]').parent().next(); nowPage.attr('style'); nowPage = nowPage.next()) {}
        nowPage = parseInt(nowPage.attr('name'))
        if (nowPage == 10) {
            show(10)
        } else {
            show(nowPage + 1)
        }
    })
})

function memberBar(name) {
    $('#memberInfo').children().remove()
    $("#memberInfo").append($('<a>').attr('class', 'nav-link').text("hi, " + name))
    $("#memberInfo").append($('<a>').attr('class', 'nav-link').attr('id', "toLogout").attr('href', "#").text("登出"))
    $("#memberInfo").append($('<a>').attr('class', 'nav-link').attr('id', "toCreateroomPage").attr('href', "#").text("建立聊天室"))

    // logout
    $("#toLogout").click(function() {
        $.ajax({
            type: "GET",
            url: "/logout",
            success: function(data) {
                location.href = "/"
            }
        });
    })

    // create room 
    $("#toCreateroomPage").click(function() {
        location.href = "/create"
    })
}
// room list
function roomTable(page) {
    for (let i = page * 10 - 9; i < page * 10 + 1; i++) {
        let td1 = $('<td>').text(i);
        let td2 = $('<td>').attr("class", "roomName")
        let a2 = $('<a>').attr('href', '/room/' + i).attr("class", "roomBtn").attr('name', i)
        let td3 = $('<td>')
        let td4 = $('<td>').attr("class", "roomEntry")
        let a4 = $('<a>').attr('href', '/room/' + i).attr("class", "roomBtn").attr('name', i)
        td2.append(a2)
        td4.append(a4)
        $('#chatroomList').append($('#' + i).append(td1, td2, td3, td4))
    }
}

function show(page) {
    let wanted = { "page": page }
    $.ajax({
        method: "POST",
        url: "/",
        data: wanted,
        beforeSend: function() {
            $('#chatroomList').empty()
            for (i = page * 10 - 9; i < page * 10 + 1; i++) {
                $('#chatroomList').append($('<tr>').attr('id', i))
            }
        },
        success: function(data) {
            roomTable(page)
            let roomList = JSON.parse(data);
            roomList.forEach(room => {
                showRoom(room.id, room.title, room.nickname)
            })
            showPage(page)
        }
    })
}

function showRoom(id, title, nickname) {
    $('#' + id + ' td:nth-child(2) a').text(title)
    $('#' + id + ' td:nth-child(3)').text(nickname)
    $('#' + id + ' td:nth-child(4) a').text("進入")
}

// page
function showPage(page) {
    $('td[class="page"]').hide()
    let totalPage = 100
    for (i = page; i <= page + 4 && page + 4 <= totalPage; i++) {
        $('td[class*="page"]' + '[name="' + i + '"]').show()
    }
}

function addPage() {
    page = 10
    let tr = $('tr[name*="pageList"]')
    tr.append(addPageTd("first", null, "第一頁"))
    tr.append(addPageTd("prev", null, "前一頁"))
    for (i = 1; i < page + 1; i++) {
        tr.append(addPageTd(null, i, null))
    }
    tr.append(addPageTd("next", null, "下一頁"))
    tr.append(addPageTd("last", null, "最後一頁"))
}

function addPageTd(name, page, text) {
    if (page == null) {
        var td = $('<td>')
        var a = $('<a>').attr('class', name).attr("href", "#").text(text)

    } else {
        var td = $('<td>').attr('class', 'page').attr('name', page)
        var a = $('<a>').attr("class", "pageBtn").attr("name", page).attr('href', '#').text(page)
    }
    td.append(a)
    return td
}