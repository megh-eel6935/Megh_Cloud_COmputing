$( document ).ready(function() {

$.get( "/getmessages", function( data ) {
  sap = JSON.stringify(data);
  var length = data.data.length;
  console.log(length);
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
  
  $( "#result" ).append(sap);
});

$("#filebutton").click(function(){
        $("#uploadform").toggle();
    }); 



});

















