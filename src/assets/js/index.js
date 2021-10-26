var fileupload;
var filekey = 'udf';

function sendPost(uid) {
  //Show loading animation
  document.getElementById('uploadbutton').style.display = 'none';
  document.getElementById('loaderAnimation').style.display = '';
  worker();

  async function worker() {
    //Stop here, this Promise make function sendPost return view update without waiting CryptoJS or axios
    //This javascript trick is not proper way, but javascript doesn't help either whith asynchronies so...
    await new Promise((r) => setTimeout(r, 0));

    //Generate pass for file
    filekey = Math.random().toString(36).substr(2);

    //Encrypt the base64 whith CryptoJS in AES
    var encryptData = await CryptoJS.AES.encrypt(document.getElementById('showimg').src, filekey).toString();

    //Generate POST petition
    var axios = window.axios;
    axios.post('./upload', {
        Name: fileupload.name,
        Base64: encryptData,
        Uid: uid,
      })
      .then(function (response) {
        //Hide loader animation
        document.getElementById('loaderAnimation').style.display = 'none';
        document.getElementById('copyButton').style.display = '';
      }).catch(
        function (error) {
          alert('Session expired 👋')
          location.reload();
        }
      );
  }
}

function share(uid) {
  //Create element to enter text
  tmpObj = document.createElement('textarea');
  tmpObj.value = window.location.href + 'share?file=' + uid + '#' + filekey;
  document.body.appendChild(tmpObj);
  tmpObj.select();
  //Copy text in element
  document.execCommand('copy');
  //Delete element
  document.body.removeChild(tmpObj);

  //User feedback
  document.getElementById('copyButton').innerHTML = 'Copied';
}

function dropHandler(ev) {
  //Drop image
  ev.preventDefault();
  if (ev.dataTransfer.items) {
    if (ev.dataTransfer.items[0].kind === 'file') {
      fileupload = ev.dataTransfer.items[0].getAsFile();
      //Update View to image
      document.getElementById('dropzone').style.display = 'none';
      document.getElementById('filename').style.display = '';
      document.getElementById('filename').innerHTML = fileupload.name;
      document.getElementById('loaderAnimation').style.display = '';

      //Launch reader for uploader image
      var reader = new FileReader();
      //Get Base64
      reader.readAsDataURL(fileupload);
      //Show reader image
      reader.onload = function () {
        document.getElementById('showimg').style.display = '';
        document.getElementById('showimg').src = reader.result;
        document.getElementById('uploadbutton').style.display = '';
        document.getElementById('loaderAnimation').style.display = 'none';
      };
    }
  }
}

function dragOverHandler(ev) {
  ev.preventDefault();
}
