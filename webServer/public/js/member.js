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





})