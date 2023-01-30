window.onload = function() {
  document.getElementById("username").focus();
}

$("#signInUser").click(function() {
  $.post(
    '/signInUser',
    {
      username: $('#username').val(),
      password: $('#password').val()
    },
    function(result) {
      if (result != 1) {
        $.showAlert(result)
      } else {
       window.location.reload()
      } 
    }
  );
    return false
})

//
function signUpUser() {
  if(($('#username').val()).length<3 || ($('#password').val()).length<3) {
    $.showAlert("Username and password fields cannot be empty", true)
  } else {
    var username = $('#username').val()
    var password = $('#password').val()
    let loginData = {
      username: username,
      password: password
    }
    let fetchData = {
      method: 'POST',
      body: JSON.stringify(loginData),
      headers: new Headers({
        'Content-Type': 'application/json; charset=UTF-8'
      })
    }
    fetch('/signUpUser', fetchData)
    .then(resposne => {
      if (resposne.status == 200) {
        $.showAlert("Successfully added new user", false)
      } else {
        $.showAlert("Could not add user to the database", true)
      }
    })
  }







  //   $.post(
  //     '/signUpUser',
  //     {
  //       username: $('#username').val(),
  //       password: $('#password').val()
  //     },
  //     function(result) {
  //       var json = $.parseJson(result)
  //       if (json.statusCode != 200) {
  //         $.showAlert("Could not add user to the database", true)
  //       } else {
  //         $.showAlert("Successfully added new user " + json.username, false)
  //       } 
  //     },
  //     "json"
  //   )
  // }
}

/// Common functions

// Show error message which pops up top of the screen
$.showAlert = function(message, isError) {
  $('div#alert').html(message);
    var $alert = $('div#alert');
    if (isError) {
      $("div#alert").css("background-color", "#F00");
    } else {
      $("div#alert").css("background-color", "#006400");
    }
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