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
            }
        }
    });
    // ws
    function initWS(roomID) {
        // var socket = new WebSocket("ws://localhost:8080/room/{[0~9]+}/echo")
        var socket = new WebSocket("ws://localhost:8080/room/" + roomID + "/echo")
        socket.onopen = function() {
            alert("socket is connected1")
            $.ajax({
                url: "/room/" + roomID + "/connRoom",
                succsess: function() {
                    console.log("connectedRoom set cookie success!")
                    alert("socket is connected")
                }
            })
        };
        socket.onmessage = function(e) {

            m = JSON.parse(e.data)
            console.log(m)
            if (m.status == 0) {
                console.log(m.status)
                removeMember2List(m.nickname)

                addMsg("System: " + m.nickname + " is disconnected!")
            } else if (m.status == 1) {
                console.log(m.status)
                addMember2List(m.nickname)
                addMsg("System: " + m.nickname + " is connected!")
            } else {
                if (m.msg !== undefined) {
                    addMsg(m.nickname + ": " + m.msg)
                }
            }
        }
        socket.onclose = function() {
            console.log("socket is disconnected")
            alert("socket is disconnected1")
            $.ajax({
                url: "/room/" + roomID + "/disconnRoom",
                succsess: function() {
                    console.log("disconnected set cookie success!")
                    alert("socket is disconnected")
                }
            })

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
})