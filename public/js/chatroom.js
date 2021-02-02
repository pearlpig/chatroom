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
            if (result == null || result.ID == -1) {
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
                let result = JSON.parse(data);
                if (result.status == 0) {
                    location.href = "/"
                } else if (result.status == 1) {
                    $("#loginEmailErrMsg").show()
                        // $('span[name*="emailErrMsg"]').css({ visibility: "visible" });
                    $("#loginEmailErrMsg").text(result.msg)
                } else if (result.status == 2) {
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
    $("#signupPage").click(function() {
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
                if (result.status == 0) {
                    location.href = "/"
                } else if (result.status == 1) {
                    $("#signupEmailErrMsg").show()
                        // $('span[name*="emailErrMsg"]').css({ visibility: "visible" });
                    $("#signupEmailErrMsg").text(result.msg)
                } else if (result.status == 2) {
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
        location.href = "/createroom"
    })
    $("#chatroomPage").click(function() {
        location.href = "/chatroom"
    })
}

function getFormData($form) {
    var unindexed_array = $form.serializeArray();
    var indexed_array = {};

    $.map(unindexed_array, function(n, i) {
        indexed_array[n['name']] = n['value'];
    });

    return indexed_array;
}