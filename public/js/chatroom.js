window.onload = function() {
    addPage()
    show(1)

    // check member
    $.ajax({
        type: "GET",
        url: "/checkMember",
        beforeSend: function() {
            $("#memberIndex").hide()
        },
        success: function(data) {
            console.log(data)
            let result = JSON.parse(data)
            console.log(data)
            if (result == null || result.member_id == -1) {
                return
            }
            $("#toLoginPage").hide()
            $("#toSignupPage").hide()
            $("#memberIndex").show()
            var oMemberBar = document.getElementById('memberIndex');
            oMemberBar.insertAdjacentText("afterbegin", "hi, " + result.nickname)
        }
    });

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
        // login in
    $("#toLoginPage").click(function() {
        location.href = "/login"
    })
    $("#loginForm").submit(function(e) {

        e.preventDefault(); // avoid to execute the actual submit of the form.

        var form = $(this);
        var url = form.attr('action');
        $.ajax({
            type: "POST",
            url: url,
            data: form.serialize(), // serializes the form's elements.
            beforeSend: function() {
                $('span[name*="emailErrMsg"]').hide()
                $('span[name*="passwordErrMsg"]').hide()
            },
            success: function(data) {
                let result = JSON.parse(data);
                if (result.code == 0) {
                    location.href = "/"
                } else if (result.code == 1) {
                    $("#loginEmailErrMsg").show()
                        // $('span[name*="emailErrMsg"]').css({ visibility: "visible" });
                    $("#loginEmailErrMsg").text(result.msg)
                } else if (result.code == 2) {
                    $("#loginPasswordErrMsg").show()
                        // $('span[name*="passwordErrMsg"]').css({ visibility: "visible" });
                    $("#loginPasswordErrMsg").text(result.msg)
                } else {
                    alert("invalid status")
                }
            }
        });


    });
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

    $("#signupForm").submit(function(e) {

        e.preventDefault(); // avoid to execute the actual submit of the form.

        var form = $(this);

        console.log(form)
        var url = form.attr('action');
        $.ajax({
            type: "POST",
            url: url,
            data: form.serialize(), // serializes the form's elements.
            beforeSend: function() {
                $("#signupEmailErrMsg").hide()
                $("#signupNicknameErrMsg").hide()
            },
            success: function(data) {
                // console.log(data)
                let result = JSON.parse(data);
                console.log(result)
                if (result.code == 0) {
                    location.href = "/"
                } else if (result.code == 1) {
                    $("#signupEmailErrMsg").show()
                        // $('span[name*="emailErrMsg"]').css({ visibility: "visible" });
                    $("#signupEmailErrMsg").text(result.msg)
                } else if (result.code == 2) {
                    $("#signupNicknameErrMsg").show()
                        // $('span[name*="passwordErrMsg"]').css({ visibility: "visible" });
                    $("#signupNicknameErrMsg").text(result.msg)

                } else {
                    alert("invalid status")
                }
            }
        });


    });

    $("#createroomPage").click(function() {
        console.log("create")
        location.href = "/create"
    })

    $("#createRoomForm").submit(function(e) {

        e.preventDefault(); // avoid to execute the actual submit of the form.

        var form = $(this);

        console.log(form)
        var url = form.attr('action');
        $.ajax({
            type: "POST",
            url: url,
            data: form.serialize(), // serializes the form's elements.
            beforeSend: function() {
                $("#createRoomNameErrMsg").hide()
            },
            success: function(data) {
                console.log(data)
                let result = JSON.parse(data);
                console.log(result)
                if (result.status.code == 0) {
                    location.href = "/room/" + result.data.id
                } else if (result.status.code == 1) {
                    $("#createRoomNameErrMsg").show()
                    $("#createRoomNameErrMsg").text(result.status.msg)
                        // $('span[name*="emailErrMsg"]').css({ visibility: "visible" });
                } else {
                    alert("invalid status")
                }
            }
        });


    });
    $("#chatroomPage").click(function() {
        location.href = "/chatroom"
    })
}

function show(page) {
    console.log(page)
    let wanted = { "page": page }
    $.ajax({
        method: "POST",
        url: "/",
        data: wanted,
        beforeSend: function() {
            $('#chatroomList').empty()
        },
        success: function(data, status) {
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
    var a1 = $('<a>').attr('href', '/room/' + id).text(title)
    var td2 = $('<td>').text(nickname);
    var td3 = $('<td>').attr("class", "roomEntry")
    var a2 = $('<a>').attr('href', '/room/' + id).text('進入')
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

// function getFormData($form) {
//     var unindexed_array = $form.serializeArray();
//     var indexed_array = {};

//     $.map(unindexed_array, function(n, i) {
//         indexed_array[n['name']] = n['value'];
//     });

//     return indexed_array;
// }