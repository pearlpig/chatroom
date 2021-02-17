let needAuth = ["/create", "/enter", "/echo"]
$(function() {
    // check member
    $.ajax({
        type: "GET",
        url: "/check",
        success: function(data) {
            result = JSON.parse(data)
            console.log(result)
            if (document.location.pathname == "/") {
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
            needAuth.forEach(url => {
                if (document.location.pathname == url) {
                    if (result.member_id == -1) {
                        location.href = "/login"
                    }
                }
            })
            let roomUrl = /\/room\/[0-9]+/
            if (roomUrl.test(document.location.pathname)) {

                if (result.member_id == -1) {
                    location.href = "/login"
                }
                if (window.WebSocket == undefined) {
                    alert("THe browser doesn't support wrbsocket!")
                } else {
                    // title = (document.title)
                    roomID = parseInt(document.location.pathname.replace(/\/room\//, ""))
                    console.log(roomID)
                    ws = initWS(roomID)
                }
                $('#sendBtn').click(function() {
                    console.log("click")
                    text = $('#chatInput').val()
                    console.log(text)
                    if (text != undefined && text !== "") {
                        ws.send(text)
                    }
                    // ws.send(JSON.stringify({ msg: $('#chatInput').val() }))
                })
                $('#quit').click(function() {
                    location.href = "/"
                })
                window.onbeforeunload = function(e) {　　
                    $.ajax({
                        url: "/room/" + roomID + "/disconnRoom",
                        succsess: function(data) {
                            console.log(data)
                            console.log("disconnected set cookie success!")
                            alert("socket is disconnected")
                        }
                    })
                }
            }
        }
    });
    // ws
    function initWS(roomID) {
        // var socket = new WebSocket("ws://localhost:8080/room/{[0~9]+}/echo")
        var socket = new WebSocket("ws://localhost:8080/room/" + roomID + "/echo")
        socket.onopen = function() {
            console.log("socket is onopen")
            $.ajax({
                url: "/room/" + roomID + "/connRoom",
            })
        };
        socket.onmessage = function(e) {
            m = JSON.parse(e.data)
            console.log(m)
            if (m.status == 0) {
                if (m.msg !== undefined) {
                    console.log("sending msg")
                    addMsg(m.nickname[0] + ": " + m.msg)
                }
            } else if (m.status == 1) {
                console.log(m.status)
                addMsg("System: " + m.nickname[0] + " is connected!")
            } else if (m.status == 2) {
                console.log(m.status)
                removeMember2List(m.nickname)
                addMsg("System: " + m.nickname[0] + " is disconnected!")
            } else if (m.status == 3) {
                removeAllMember2List()
                m.nickname.forEach(name => {
                    addMember2List(name)
                })
            }
        }
        socket.onclose = function() {
            // addMsg("Socket is close", "System")
        }
        return socket
    }

    function addMember2List(n) {
        $('.roomMemberList').append($('<li>').attr('class', 'listMember').attr('name', n).text(n))
    }

    function removeMember2List(n) {
        $('li[name="' + n + '"]').remove()
    }

    function removeAllMember2List() {
        $('.roomMemberList').children().remove()
    }

    function addMsg(m) {
        console.log("add")
        let ul = $('#chatContent')
        ul.append($('<li>').text(m))
    }
    // login in
    $("#toLoginPage").click(function() {
        location.href = "/login"
    })
    $("#loginForm").submit(function(e) {

        e.preventDefault(); // avoid to execute the actual submit of the form.

        var form = $(this);
        var url = form.attr('action');
        let email, pwd
        data = form.serializeArray()
        data.forEach(item => {
            console.log(item)
            if (item.name == "email") {
                email = item.value
            } else if (item.name == "pwd") {
                pwd = item.value
            }
        })
        console.log(email, pwd)
        $.ajax({
            type: "POST",
            url: url,
            data: form.serialize(), // serializes the form's elements.
            beforeSend: function() {
                $('#loginEmailErrMsg').hide()
                $('#loginPasswordErrMsg').hide()
                if (checkEmailFmt(email) !== "ok") {
                    $('#loginEmailErrMsg').show()
                    $('#loginEmailErrMsg').text(checkEmailFmt(email))
                    return false
                } else if (checkPwdFmt(pwd) !== "ok") {
                    $('#loginPasswordErrMsg').show()
                    $('#loginPasswordErrMsg').text(checkPwdFmt(pwd))
                    return false
                }
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

        var url = form.attr('action');
        let email, nickname, pwd1, pwd2
        data = form.serializeArray()
        data.forEach(item => {
            console.log(item)
            if (item.name == "email") {
                email = item.value
            } else if (item.name == "nickname") {
                nickname = item.value
            } else if (item.name == "pwd1") {
                pwd1 = item.value
            } else if (item.name == "pwd2") {
                pwd2 = item.value
            }
        })

        $.ajax({
            type: "POST",
            url: url,
            data: form.serialize(), // serializes the form's elements.
            beforeSend: function() {
                $("#signupEmailErrMsg").hide()
                $("#signupNicknameErrMsg").hide()
                $("#passwordErrMsg").hide()
                $("#confirmdPasswordErrMsg").hide()
                if (checkEmailFmt(email) !== "ok") {
                    console.log(checkEmailFmt(email))
                    $("#signupEmailErrMsg").show()
                    $("#signupEmailErrMsg").text(checkEmailFmt(email))
                    return false
                } else if (checkNicknameFmt(nickname) !== "ok") {
                    console.log(checkNicknameFmt(nickname))
                    $("#signupNicknameErrMsg").show()
                    $("#signupNicknameErrMsg").text(checkNicknameFmt(nickname))
                    return false
                } else if (checkPwdFmt(pwd1) !== "ok") {
                    $("#passwordErrMsg").show()
                    $("#passwordErrMsg").text(checkPwdFmt(pwd1))
                    return false
                } else if (checkPwdConfirm(pwd1, pwd2) !== "ok") {
                    $("#confirmdPasswordErrMsg").show()
                    $("#confirmdPasswordErrMsg").text(checkPwdConfirm(pwd1, pwd2))
                    return false
                }
            },
            success: function(data) {
                console.log(data)
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
})

function checkEmailFmt(email) {
    return "ok"
}

function checkNicknameFmt(nickname) {
    if (nickname.length > 20) {
        return "Nickname length should at most 20 character!"
    } else if (nickname.length < 1) {
        return "Nickname should not be empty!"
    }
    return "ok"
}

function checkPwdFmt(pwd) {

    if (pwd.length < 8 || pwd.length > 128) {
        return "Password length should between 8 to 128!"
    }
    return "ok"
}

function checkPwdConfirm(pwd1, pwd2) {
    if (pwd1 !== pwd2) {
        return "Please check the confirmed password!"
    }
    return "ok"
}