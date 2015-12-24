username = Bacon.UI.textFieldValue($("#username"))
email = Bacon.UI.textFieldValue($("#email"))

function notEmpty(x) { return x.length > 0 }

usernameEntered = username.map(notEmpty)
emailEntered = email.map(notEmpty)

buttonEnabled = usernameEntered.and(emailEntered)
buttonEnabled.onValue(function(enabled) {
    $("#submit").attr("disabled", !enabled)
})
