var ChatRoom = function() {
  var c = this;
  
  c.create = function() {
    c.sendForm   = document.getElementById("send-form");
    c.textInput  = c.sendForm.querySelector("#groupmessage");
    
    return c;
  };
  
  c.connect = function() {
    


console.log("ws://" + location.host + "/sockets/sanath");
    c.ws = new WebSocket("ws://" + location.host + "/sockets/sanath");


    return c;
  };
  
  c.disconnect = function() {
    c.ws.close();
  };
  
  c.listen = function() {
    c.ws.addEventListener("message", function(e) {
      console.log("received"+e.data);
      
      var temp=JSON.parse(e.data);
      $("#log").append("<br><div align='left'>"+temp.text+"<br>");
      updatelist(temp.groupid);

    });
    
    c.ws.addEventListener("close", function(e) {
      console.log("closed connection, reconnecting", e)
      c.connect()
    });
    
    c.ws.addEventListener("error", function(e) {
      console.log("error for connection", e)
    });
    
    return c;
  };
  
   c.send = function(messageData) {
    console.log("sent"+JSON.stringify(messageData));
     $("#log").append("<br><div align='right'>"+messageData.text+"</div>");   
    c.ws.send(JSON.stringify(messageData));
  };

   c.onopen = function() {
      //c.ws.send($.toJSON({ action: 'join', user: name }));
      console.log("socket opened")
    }
  
  c.handleSend = function() {
    c.sendForm.addEventListener("click", function(e) {
      e.preventDefault();
      
      if (c.textInput.value && c.textInput.value.length) {
        var messageData = {
          type: "message",
          from : "sanath",
          groupid:$("#addgroups").val(),
          text: c.textInput.value
        };
        
        c.send(messageData);
        
        
        c.textInput.value = ''
      }
      
      return false;
    });
    
    return c;
  };
  

  var initialize = function() {

    if(c.ws==undefined){
   // Do your stuff...
      console.log("c.ws value "+c);
      c.create().handleSend();
      c.connect().listen();
      console.log(c.ws.readyState);
      }
      else{

          console.log(c.ws.readyState);
      }
   
  }();
};

$( document ).ready(function() {
chat = new ChatRoom();

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
      //location.reload();
      $("#groupmessage").value="";
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
       /* var response = jQuery.parseJSON(data);
        if(response.code == "success") {
            alert("Success!");
        } else if(response.code == "failure") {
            alert(response.err);
        }
        */
    }
});



});



$("#filebutton").click(function(){
        $("#uploadform").toggle();
    }); 

});

function updatelist(id){
  var url = "/getgroupdatabyid/"+id;
$.get( url, function( data ) {
  sap = JSON.stringify(data);
  var length = data.data.length;
  console.log("from update list"+length);
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
}