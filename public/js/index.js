$(function() {
    addPage()
    show(1)
    $.ajax({
        type: "GET",
        url: "/check",
        success: function(data) {
            result = JSON.parse(data)
            console.log(result)
            $("#memberInfo").hide()
            if (result.member_id != -1) {
                $("#toLoginPage").hide()
                $("#toSignupPage").hide()
                $("#memberInfo").show()
                var oMemberBar = document.getElementById('memberInfo');
                if (oMemberBar != null) {
                    oMemberBar.insertAdjacentText("afterbegin", "hi, " + result.nickname)
                }
            }

        }
    });
    $("#createroomPage").click(function() {
        location.href = "/create"
    })
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
        for (nowPage = $('a[class*="prev"]').parent().next(); nowPage.attr('style'); nowPage = nowPage.next()) {

        }
        nowPage = parseInt(nowPage.attr('name'))
        if (nowPage == 10) {
            show(10)

        } else {
            show(nowPage + 1)
        }
    })
    $("#toLoginPage").click(function() {
        location.href = "/login"
    })

    // logout
    $("#logout").click(function() {
            $.ajax({
                type: "GET",
                url: "/logout",
                success: function(data) {
                    location.href = "/"
                }
            });
        })
        // sign up 
    $("#toSignupPage").click(function() {
        console.log("signup")
        location.href = "/signup"
    })
})

function show(page) {
    let wanted = { "page": page }
    $.ajax({
        method: "POST",
        url: "/",
        data: wanted,
        beforeSend: function() {
            $('#chatroomList').empty()
        },
        success: function(data) {
            let roomList = JSON.parse(data);
            roomList.forEach(room => {
                showRoom(room.id, room.title, room.nickname)
            })
            showPage(page)
        }
    })
}


function showRoom(id, title, nickname) {
    var td1 = $('<td>').attr("class", "roomName")
    var a1 = $('<a>').attr('href', '/room/' + id).attr("class", "roomBtn").attr('name', id).text(title)
    var td2 = $('<td>').text(nickname);
    var td3 = $('<td>').attr("class", "roomEntry")
    var a2 = $('<a>').attr('href', '/room/' + id).attr("class", "roomBtn").attr('name', id).text('進入')
    td1.append(a1)
    td3.append(a2)
    var tr = $('<tr>').attr("class", "room").append(td1, td2, td3);
    $('#chatroomList').append(tr);
}

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