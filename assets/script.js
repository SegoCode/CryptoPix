var fileupload;



function encryptAndSend(uid) {

}

function sendPost(data, uid) {
	var axios = window.axios;
	axios.post('./upload', {
		Name: fileupload.name,
		Base64: data,
		Uid: uid
	  }).then(function (response) {
		console.log(response);
	  })

}




function dropHandler(ev) {
  ev.preventDefault();
	  if (ev.dataTransfer.items) {
		if (ev.dataTransfer.items[0].kind === 'file') {
		  fileupload = ev.dataTransfer.items[0].getAsFile();		   
		   //Update View
		   document.getElementById("dropzone").style.display = "none";
		   document.getElementById("filename").style.display = "";
		   document.getElementById('filename').innerHTML = fileupload.name;
           var reader = new FileReader();
		   reader.readAsDataURL(fileupload);
			reader.onload = function () {
				document.getElementById("showimg").style.display = "";
				document.getElementById("showimg").src = reader.result; 
				document.getElementById("uploadbutton").style.display = "";
			};
		}

	  } 
}


function dragOverHandler(ev) {
  ev.preventDefault();
  console.log('File(s) in drop zone');
}


