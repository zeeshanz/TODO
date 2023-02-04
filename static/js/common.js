// All pages have body with opacity 0 and they are faded in for a nice effect.
window.onload = function () {
  disableButtonIfFieldsAreEmpty()
  $('#container').fadeTo("slow", 1);
  if($('#username').length) { // valid only for index.html
    $('#username').focus
  }
}

//
function signInUser() {
  if (($('#username').val()).length < 3 || ($('#password').val()).length < 3) {
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
          var alerttimer = window.setTimeout(function () {
            $("html").fadeOut(function () {
              window.location.href = "/tasks";
            });
          }, 2000);
        } else {
          $.showAlert("Sign in failed. Username or password incorrect.", true)
        }
      })
  }
}

//
function signUpUser() {
  if (($('#username').val()).length < 3 || ($('#password').val()).length < 3) {
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

//
function signOutUser() {
  $.showAlert("Signing out", false)
  $('#container').fadeOut
  var alerttimer = window.setTimeout(function () {
    $("html").fadeOut(function () {
      window.location.href = "/signOutUser";
    });
  }, 1000);
}

//
function addTodoItem() {
  var todoItem = document.getElementById("todoItem").value;
  fetch('/addNewTodo', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ completed: false, task_name: todoItem })
  }).then(resposne => {
    console.log(resposne)
  })
}

/// COMMON FUNCTIONS

// Show error message which pops up top of the screen
$.showAlert = function (message, isError) {
  $('div#alert').html(message);
  var $alert = $('div#alert');
  if (isError) {
    $("div#alert").css("background-color", "#F00");
  } else {
    $("div#alert").css("background-color", "#006400");
  }
  if ($alert.length) {
    var alerttimer = window.setTimeout(function () {
      $alert.trigger('click');
    }, 2000);
    $alert.animate({ height: $alert.css('line-height') || '50px' }, 200)
      .click(function () {
        window.clearTimeout(alerttimer);
        $alert.animate({ height: '0' }, 200);
      });
  }
}

// Enable or disable the sign in button if input fields are empty
function disableButtonIfFieldsAreEmpty() {
  $('.forminput').keyup(function () {
    var empty = false;
    $('.forminput').each(function () {
      if ($(this).val().length < 4) {
        empty = true;
      }
    });

    if (empty)
      $('#blueButton').prop('disabled', true);
    else
      $('#blueButton').prop('disabled', false);
  });
}