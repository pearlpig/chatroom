// $(function() {
//     let roomUrl = /\/room\/[0-9]+/
//     if (roomUrl.test(document.location.pathname)) {
//         if (window.WebSocket == undefined) {
//             alert("THe browser doesn't support wrbsocket!")
//         } else {
//             // title = (document.title)
//             ws = initWS(1)
//         }
//         $('#sendBtn').click(function() {
//             console.log("click")
//             text = $('#chatInput').val()
//             console.log(text)
//             if (text != undefined && text !== "") {
//                 ws.send(text)
//             }
//             // ws.send(JSON.stringify({ msg: $('#chatInput').val() }))
//         })
//         $('#quit').click(function() {
//             location.href = "/"
//         })
//     }


// })

// function initWS(roomID) {
//     // var socket = new WebSocket("ws://localhost:8080/room/{[0~9]+}/echo")
//     var socket = new WebSocket("ws://localhost:8080/room/" + roomID + "/echo")
//     socket.onopen = function() {
//         $.ajax({
//             data: { roomID: roomID },
//             url: "/connectdRoom",
//             succsess: function() {
//                 location.href("/")
//             }
//         })
//     };
//     socket.onmessage = function(e) {

//         m = JSON.parse(e.data)
//         console.log(m)
//         if (m.status == 0) {
//             console.log(m.status)
//             removeMember2List(m.nickname)

//             addMsg("System: " + m.nickname + " is disconnected!")
//         } else if (m.status == 1) {
//             console.log(m.status)
//             addMember2List(m.nickname)
//             addMsg("System: " + m.nickname + " is connected!")
//         } else {
//             if (m.msg !== undefined) {
//                 addMsg(m.nickname + ": " + m.msg)
//             }
//         }
//     }
//     socket.onclose = function() {
//         console.log("socket is disconnected")
//         $.ajax({
//             data: { roomID: roomID },
//             url: "/disconnectdRoom",
//             succsess: function() {
//                 location.href("/")
//             }
//         })

//         // addMsg("Socket is close", "System")
//     }
//     return socket
// }

// function addMember2List(n) {
//     $('.roomMemberList').append($('<li>').attr('class', 'listMember').attr('name', n).text(n))
// }

// function removeMember2List(n) {
//     $('li[name="' + n + '"]').remove()
// }

// function addMsg(m) {
//     console.log("add")
//     let ul = $('#chatContent')
//     ul.append($('<li>').text(m))
// }