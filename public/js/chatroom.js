window.onload = function() {
    // check member
    $.ajax({
        type: "GET",
        url: "/checkMember",
        beforeSend: function() {
            $("#memberIndex").hide()
        },
        success: function(data) {
            console.log(data)
            let result = JSON.parse(data);
            if (result.ID == -1) {
                return
            }
            console.log(result)
            $("#loginPage").hide()
            $("#signupPage").hide()
            $("#memberIndex").show()
            var oMemberBar = document.getElementById('memberIndex');
            oMemberBar.insertAdjacentText("afterbegin", "hi, " + result.Nickname)
        }
    });
    $('span[name*="emailErrMsg"]').hide()
    $('span[name*="passwordErrMsg"]').hide()

    // login in
    $("#loginPage").click(function() {
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
                console.log(data)
                let result = JSON.parse(data);
                console.log(result)
                if (result.status == 0) {
                    location.href = "/"
                } else if (result.status == 1) {
                    $('span[name*="emailErrMsg"]').show()
                        // $('span[name*="emailErrMsg"]').css({ visibility: "visible" });
                    $('span[name*="email"]').text(result.msg)
                } else if (result.status == 2) {
                    $('span[name*="passwordErrMsg"]').show()
                        // $('span[name*="passwordErrMsg"]').css({ visibility: "visible" });
                    $('span[name*="password"]').text(result.msg)

                } else {
                    alert("invalid status")
                }
            }
        });


    });

    $("#logout").click(function() {
        $.ajax({
            type: "GET",
            url: "/logout",
            success: function(data) {
                location.href = "/"
            }
        });
    })
    $("#signupPage").click(function() {
        location.href = "/signup"
    })
    $("#createroomPage").click(function() {
        location.href = "/createroom"
    })
    $("#chatroomPage").click(function() {
        location.href = "/chatroom"
    })
}