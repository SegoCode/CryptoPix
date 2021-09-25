function loadView() {
  if(window.location.hash) {
	  //Get url fragment
      var hash = window.location.hash.substring(1); 
	  
	  //Decrypt from data in src img element
	  var imgDecrypted = CryptoJS.AES.decrypt(document.getElementById("showimg").src, hash);
	  
	  //Rewrite src
	  document.getElementById("showimg").src = imgDecrypted.toString(CryptoJS.enc.Utf8);
	  document.getElementById("loaderAnimation").style.display = "none";
  } else {
	  alert("The url not contains decryption password :(")
	  location.href = '/';
  }
}