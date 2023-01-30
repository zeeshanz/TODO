window.onload = function() {
  document.getElementById("username").focus();
}

$("#signUpUser").click(function() {
  $.post(
    '/actions/login.php',
    {
      username: $('#username').val(),
      password: $('#password').val()
    },
    function(result) {
      if (result != 1) {
        $.terror(result)
      } else {
       window.location.reload()
      } 
    }
  );
    return false
})


function signUpUser() {
  if(($('#username').val()).length<3 || ($('#password').val()).length<3) {
    $.terror("Username and password fields cannot be empty")
  } else {
    $('#signUpUser').signUpUser()
  }
}

/// Common functions
$.terror = function(message) {
  $('div#error').html(message);
    var $alert = $('div#error');
    if($alert.length) {
      var alerttimer = window.setTimeout(function () {
          $alert.trigger('click');
      }, 3000);
      $alert.animate({height: $alert.css('line-height') || '50px'}, 200)
      .click(function () {
          window.clearTimeout(alerttimer);
          $alert.animate({height: '0'}, 200);
      });
  }
}