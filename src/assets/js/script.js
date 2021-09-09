var fileupload;

function sendPost(uid) {
	
	//var encryptImage = CryptoJS.AES.encrypt(document.getElementById("showimg").src, 'secret key');
	document.getElementById("uploadbutton").style.display = "none";
	document.getElementById("loaderAnimation").style.display = "";
	
	(async () => {
		var axios = window.axios;
		axios.post('./upload', {
			Name: fileupload.name,
			Base64: document.getElementById("showimg").src,
			Uid: uid
		  }).then(function (response) {
			document.getElementById("loaderAnimation").style.display = "none";
			document.getElementById("uploadbutton").style.display = "";
			document.getElementById("uploadbutton").innerHTML = "Copy link";
		  })
	})();
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


