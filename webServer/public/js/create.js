$(function() {
    $.ajax({
        type: "GET",
        url: "/check",
        success: function(data) {
            console.log("cehckout")
            result = JSON.parse(data)
            console.log(result)
            if (result.member_id == -1) {
                location.href = "/login"
            }
        }
    });
    $("#createRoomForm").submit(function(e) {
        e.preventDefault(); // avoid to execute the actual submit of the form.
        var form = $(this);
        var url = form.attr('action');
        formData = getFormData(form)
        $.ajax({
            type: "POST",
            url: url,
            data: form.serialize(), // serializes the form's elements.
            beforeSend: function() {
                if (checkRoomNameFmt(formData.roomName) !== "ok") {
                    errMsg(checkRoomNameFmt(formData.roomName))
                    return false
                }
            },
            success: function(data) {
                console.log(data)
                let result = JSON.parse(data);
                console.log(result)
                if (result.status.code == 0) {
                    location.href = "/"
                } else if (result.status.code == 1) {
                    $("#createRoomNameErrMsg").show()
                    $("#createRoomNameErrMsg").text(result.status.msg)
                } else {
                    alert("invalid status")
                }
            }
        });

    });
})

function checkRoomNameFmt(roomName) {

    if (roomName.length > 20) {
        return "Room name length should at most 20 character."
    } else if (roomName.length < 1) {
        return "Room name should not be empty!"
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