$( document ).ready(function() {

$("#addgroups").click(function(){
  var id =$("#addgroups").val();
  var url = "/getgroupdatabyid/"+id;
$.get( url, function( data ) {
  sap = JSON.stringify(data);
  var length = data.data.length;
  console.log(length);
  $("#k").empty();
  for(var i=0;i<length;i++)
  {
    if(data.data[i].filetype=="text"){
  	     var temp="<tr><td>"+data.data[i].message+" <td> <td align='right'>message</td></tr>";
       }
    else if(data.data[i].filetype=="file") {
      var temp="<tr><td>"+data.data[i].message+"<td> <td align='right'><a href='getfile?filename="+data.data[i].message+"'>download</a></td></tr>";
    }
  	$("#k").append(temp);
  }
  
});
});


$.get( "/groupslist", function( data ) {
  sap = JSON.stringify(data);
  var length = data.data.length;
  ;
  for(var i=0;i<length;i++)
  {
    
    var temp="<li><button type=\"button\" class=\"btn btn-default\" value=\""+data.data[i].groupid+"\">"+data.data[i].groupname+"</button></li>";
    var temp = "<option value="+data.data[i].groupid+">"+data.data[i].groupname+"</option>";
    $("#addgroups").append(temp);
  }
  id =$("#addgroups").val();
  url = "/groupdata/"+id;

});
  


$("#groupsubmit").click(function(){
  var value =$("#groupmessage").val();
  var id =$("#addgroups").val();
  var url = "/groupdata/"+id;
  $.post(url,{message:value},function(data,status){
      location.reload();
  });
});

$("#groupfilesubmit").click(function(){
  var value =$("#groupmessage").val();
  var id =$("#addgroups").val();
  var url = "/groupuploadfile/"+id;
  var file = $("#uploadfile")[0].files[0];
var formdata = new FormData();
formdata.append("file", file);
 $.ajax({
  url: url,
  data: formdata,
  processData: false,
  contentType: false,
  type: 'POST',
  success: function(data) {
      
    }
});


});



$("#filebutton").click(function(){
        $("#uploadform").toggle();
    }); 

});
