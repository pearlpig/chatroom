$(function() {
    $("#loginForm").submit(function(e) {
        e.preventDefault(); // avoid to execute the actual submit of the form.

        var form = $(this);
        formData = getFormData(form)
        $.ajax({
            type: "POST",
            url: "/login",
            data: form.serialize(), // serializes the form's elements.
            beforeSend: function() {
                if (checkEmailFmt(formData.email) !== "ok") {
                    errMsg(checkEmailFmt(formData.email))
                    return false
                } else if (checkPwdFmt(formData.pwd) !== "ok") {
                    errMsg(checkPwdFmt(formData.pwd))
                    return false
                }
            },
            success: function(data) {
                console.log(data)
                let result = JSON.parse(data);
                if (result.code == 0) {
                    location.href = "/"
                } else {
                    errMsg(result.msg)
                }
            }
        });
    });
})

function checkEmailFmt(email) {
    return "ok"
}

function checkPwdFmt(pwd) {
    if (pwd.length < 8) {
        return "Password length should at least 8."
    }
    return "ok"
}

function errMsg(err) {
    $('#errBox').children().remove()
    $('#errBox').append('<span>').attr('style', "color:red; font-size: medium;").text(err)
}

function getFormData($form) {
    var unindexed_array = $form.serializeArray();
    var indexed_array = {};

    $.map(unindexed_array, function(n, i) {
        indexed_array[n['name']] = n['value'];
    });

    return indexed_array;
}