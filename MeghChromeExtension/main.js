//chrome.browserAction.onClicked.addListener(function (t) {
//window.onload = function () {
console.log("in the file!!!!!");
var button = document.getElementById("upload_button");
$(button).click(function(){
    var formData = new FormData($('form')[0]);
    $.ajax({
  url: 'http://104.131.126.89/uploadfile', 
  type: 'POST',
  data: formData, // The form with the file inputs.
  processData: false ,                         // Using FormData, no need to process data.
  success: function(response){
console.log("success");
}
  /*error: function(response){
console.log("upload failed");
}*/
});
});
//});