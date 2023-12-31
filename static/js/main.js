var Search = function() {
  $("#warning").text("").removeClass("visible").addClass("hidden");
  $("#submit").addClass('disabled');
  var search_word = $('#search_word').val();
  if (!search_word) {
    $("#submit").removeClass('disabled');
    $("#warning").text("Message is Empty").removeClass("hidden").addClass("visible");
    return false;
  }
  const data = {action: "search", word: search_word};
  request(data, (res)=>{
    $("#info").removeClass("hidden").addClass("visible");
    if (!!res.list) {
      let tmp = "";
      res.list.forEach(v => {
        tmp += '<a class="item">' + v + '</a>'
      });
      $("#result").html(tmp);
    }
    $("#submit").removeClass('disabled');
  }, (e)=>{
    console.log(e.responseJSON.message);
    $("#warning").text(e.responseJSON.message).removeClass("hidden").addClass("visible");
    $("#submit").removeClass('disabled');
  });
};

var request = function(data, callback, onerror) {
  $.ajax({
    type:          'POST',
    dataType:      'json',
    contentType:   'application/json',
    scriptCharset: 'utf-8',
    data:          JSON.stringify(data),
    url:           App.url
  })
  .done(function(res) {
    callback(res);
  })
  .fail(function(e) {
    onerror(e);
  });
};
var App = { url: location.origin + {{ .ApiPath }} };
