$(function() {
    $("#signupForm").submit(function(e) {
        console.log("signup")
        e.preventDefault(); // avoid to execute the actual submit of the form.
        var form = $(this);
        var url = form.attr('action');
        formData = getFormData(form)
        $.ajax({
            type: "POST",
            url: url,
            data: form.serialize(), // serializes the form's elements.
            beforeSend: function() {
                $("#signupEmailErrMsg").hide()
                $("#signupNicknameErrMsg").hide()
                $("#passwordErrMsg").hide()
                $("#confirmdPasswordErrMsg").hide()
                if (checkEmailFmt(formData.email) !== "ok") {
                    console.log(checkEmailFmt(email))
                    $("#signupEmailErrMsg").show()
                    $("#signupEmailErrMsg").text(checkEmailFmt(formData.email))
                    return false
                } else if (checkNicknameFmt(nickname) !== "ok") {
                    console.log(checkNicknameFmt(nickname))
                    $("#signupNicknameErrMsg").show()
                    $("#signupNicknameErrMsg").text(checkNicknameFmt(formData.nickname))
                    return false
                } else if (checkPwdFmt(formData.pwd1) !== "ok") {
                    $("#passwordErrMsg").show()
                    $("#passwordErrMsg").text(checkPwdFmt(formData.pwd1))
                    return false
                } else if (checkPwdConfirm(formData.pwd1, formData.pwd2) !== "ok") {
                    $("#confirmdPasswordErrMsg").show()
                    $("#confirmdPasswordErrMsg").text(checkPwdConfirm(formData.pwd1, formData.pwd2))
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

function getFormData($form) {
    var unindexed_array = $form.serializeArray();
    var indexed_array = {};

    $.map(unindexed_array, function(n, i) {
        indexed_array[n['name']] = n['value'];
    });

    return indexed_array;
}