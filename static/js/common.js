// All pages have body with opacity 0 and they are faded in for a nice effect.
window.onload = function () {
  disableButtonIfFieldsAreEmpty()
  $('#container').fadeTo("slow", 1)
  if ($('#username').length) { // valid only for index.html
    $('#username').focus
  }

  // valid only for todos.html to strike through the completed todos
  if ($('#todoItems').length) {
    $('#todoItems > tbody > tr').each(function (data) {
      var $this = $(this);
      var completed = $this.data('completed'); // or var filter = $this.attr('data-filter')
      var uuid = $(this).attr("id")
      setCompleted(uuid, completed)
    })
  }
}

// Authenticate and if successful sign in the user, otherwise show error message
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
      .then(response => {
        var status = response.status
        if (status == 200) {
          $.showAlert("Sign in successful. Opening the ToDo page", false)
          var alerttimer = window.setTimeout(function () {
            $("html").fadeOut(function () {
              window.location.href = "/todos"
            })
          }, 2000)
        } else {
          $.showAlert("Sign in failed. Username or password incorrect.", true)
        }
      })
  }
}

// Sign up a new user and show success or failure message
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
      .then(response => response.json())
      .then(response => {
        var json = JSON.parse(JSON.stringify(response))
        $.showAlert(json.message, !response.success)
      })
  }
}

// Signout the user and go back to the sign in screen
function signOutUser() {
  $.showAlert("Signing out", false)
  $('#container').fadeOut
  var alerttimer = window.setTimeout(function () {
    $("html").fadeOut(function () {
      window.location.href = "/signOutUser"
    })
  }, 1000)
}

// Add a new Todo item and append to existing Todos with a nice animation
function addTodoItem() {
  var todoItem = document.getElementById("todoItem").value
  fetch('/addNewTodo', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ completed: false, todo_item: todoItem })
  }).then(response => response.json())
    .then(response => {
      if (response.status == 200) {
        var json = JSON.parse(JSON.stringify(response))
        var uuid = json.uuid
        var todoItem = json.todoItem
        var newRow = $("<tr id='" + uuid + "'><td align='left'><a id='" + uuid + "' onclick='editTodo(id)'> ??? </a><a id='" + uuid + "' onclick='completeTodo(id)'> ??? </a><a id='" + uuid + "' onclick='deleteTodo(id)'> ??? </a><input class='inputdisabled' id='span" + uuid + "' value='" + todoItem + "' disabled='true'/><a hidden id='save" + uuid + "' onclick='updateTodo(id)'> ???? </a></td></tr>")
        newRow.hide()
        $('#todoItems tr').last().after(newRow)
        newRow.fadeIn("slow")
        $('#todoItem').val('')
        $('#blueButton').prop('disabled', true)
      } else {
        $.showAlert("An error has occured: Error code: " + response.status, true)
      }
    })
}

// Delete a Todo item and update the UI
function deleteTodo(uuid) {
  fetch('/deleteTodo', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ uuid: uuid })
  }).then(response => {
    if (response.status == 200) {
      $('#' + uuid).fadeTo("slow", 0.0, function () {
        $(this).remove()
      })
    } else {
      $.showAlert("An error has occured: Error code: " + response.status, true)
    }
  })
}

// Mark Todo as completed or not completed and update the UI
function completeTodo(uuid) {
  fetch('/completeTodo', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ uuid: uuid })
  }).then(response => {
    switch (response.status) {
      case 201:
        setCompleted(uuid, false)
        break;
      case 202:
        setCompleted(uuid, true)
        break;
      default:
        $.showAlert("An error has occured: Error code: " + response.status, true)
    }
  })
}

// Chage the UI elements to allow user to update the Todo item
function editTodo(uuid) {
  if ($('#span' + uuid).prop('disabled')) {
    $('#span' + uuid).prop('disabled', false).removeClass("inputdisabled").addClass("inputenabled").focus()
    $('#save' + uuid).show()
  } else {
    $('#span' + uuid).prop('disabled', true).removeClass("inputenabled").addClass("inputdisabled")
    $('#save' + uuid).hide()
  }
}

// Send to the server updated todo item text. If text size is less then 4 characters an error will be returned
function updateTodo(uuid) {
  var inputId = uuid.replace('save', '')
  var newTodo = $('#span' + inputId).val()
  if (newTodo.length < 4) {
    $.showAlert("Todo item must be at least 4 characters long", true)
  } else {
    fetch('/updateTodo', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ uuid: inputId, todo_item: newTodo })
    }).then(response => {
      if (response.status == 200) {
        $.showAlert("Todo updated", false)
        $('#span' + inputId).prop('disabled', true).removeClass("inputenabled").addClass("inputdisabled")
        $('#save' + inputId).hide()
      } else {
        if (response.status == 403) {
          $.showAlert("Length of text too short. Error code: " + response.status, true)
        } else {
          $.showAlert("An error occurred. Error code: " + response.status, true)
        }
      }
    })
  }
}

/// COMMON FUNCTIONS

// Show error message which pops up top of the screen
$.showAlert = function (message, isError) {
  $('div#alert').html(message)
  var $alert = $('div#alert')
  if (isError) {
    $("div#alert").css("background-color", "#F00")
  } else {
    $("div#alert").css("background-color", "#006400")
  }
  if ($alert.length) {
    var alerttimer = window.setTimeout(function () {
      $alert.trigger('click')
    }, 2000)
    $alert.animate({ height: $alert.css('line-height') || '50px' }, 200)
      .click(function () {
        window.clearTimeout(alerttimer)
        $alert.animate({ height: '0' }, 200)
      })
  }
}

// Enable or disable the sign in button if input fields are empty
function disableButtonIfFieldsAreEmpty() {
  $('.forminput').keyup(function () {
    var empty = false
    $('.forminput').each(function () {
      if ($(this).val().length < 4) {
        empty = true
      }
    })

    if (empty)
      $('#blueButton').prop('disabled', true)
    else
      $('#blueButton').prop('disabled', false)
  })
}

// Add or remove line through based on whether the todo item is completed or not
function setCompleted(uuid, completed) {
  if (completed == true) {
    $('#span' + uuid).css({ textDecoration: 'line-through' })
  } else {
    $('#span' + uuid).css({ textDecoration: 'none' })
  }
}