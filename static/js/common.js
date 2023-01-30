window.onload = function() {
  document.getElementById("username").focus();
  disableButtonIfFieldsAreEmpty()
}

function signInUser() {
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
    fetch('/signInUser', fetchData)
    .then(resposne => {
      var status = resposne.status
      if (status == 200) {
        $.showAlert("Sign in successful. Opening the ToDo page", false)
      } else {
        $.showAlert("Sign in failed", true)
      }
    })
  }
}

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
      var status = resposne.status
      if (status == 200) {
        $.showAlert("Successfully added new user", false)
      } else {
        $.showAlert("Could not add user to the database. Response code: " + status, true)
      }
    })
  }
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

//
function disableButtonIfFieldsAreEmpty() {
  $('.forminput').keyup(function () {
    var empty = false;
    $('.forminput').each(function () {
      if ($(this).val().length < 4) {
        empty = true;
      }
    });

    if (empty)
      $('#signInUserButton').prop('disabled', true);
    else
      $('#signInUserButton').prop('disabled', false);
  });
}