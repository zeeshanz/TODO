<script>
    
    window.onload = function() {
    document.getElementById("username").focus()
    }
    
    $("#usernamesubmit").click(function() {
      $.post(
        '/actions/login.php',
        {
          username: $('#username').val(),
          password: $('#password').val()
        },
        function(result) {
          if (result != 1) {
            $.terror(result);
          } else {
           window.location.reload()
          }
        }
      )
        return false;
    })
    
    $.terror = function(message) {
    $('div#error').html(message)
      var $alert = $('div#error')
      if($alert.length) {
        var alerttimer = window.setTimeout(function () {
            $alert.trigger('click');
        }, 3000)
        $alert.animate({height: $alert.css('line-height') || '50px'}, 200)
        .click(function () {
            window.clearTimeout(alerttimer);
            $alert.animate({height: '0'}, 200);
        })
    }
  }
  </script>