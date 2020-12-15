var modal = document.getElementById("myModal");
var btn = document.getElementsByName("myBtn");
var img = document.getElementById("myImg");

var imgLayer = document.getElementById("modalImage");
var selectLayer = document.getElementById("modalSelect");
var configLayer = document.getElementById("modalConfig");

var timer;
var istrue = false;
var delay = 1000;
var isCheckMode = false;
var checkedCount = 0;
var imgMode = false;

function thumbClick(id)
{
  var element = document.getElementById(id);  

  if(isCheckMode)
  {
    if(element.classList.contains("checked_img"))
    {
      element.classList.remove("checked_img");
      checkedCount--;
    }
    else
    {
      element.classList.add("checked_img");
      checkedCount++;
    }
    if(checkedCount == 0)
      DisableCheckMode();
  }
  else
  {
    modal.style.display = "block";
    imgLayer.style.display = "block";
    selectLayer.style.display = "none";
    configLayer.style.display = "none";
    
    img.src = element.name;
    img.name = id;
    imgMode = true;
  }

  console.log('count : ' + checkedCount);
}

function setting()
{
  modal.style.display = "block";
  imgLayer.style.display = "none";
  selectLayer.style.display = "none";
  configLayer.style.display = "block";
}

function movdImage()
{
  modal.style.display = "block";
  imgLayer.style.display = "none";
  selectLayer.style.display = "block";
  configLayer.style.display = "none";
}

window.onclick = function(event) {
  if (event.target == modal) {
    modal.style.display = "none";
    imgMode = false;
  }
}

function func(id)
{
   istrue = true;
   timer = setTimeout(function(){ makeChange(id);},delay);
}

function makeChange(id)
{
  if(timer)
    clearTimeout(timer);
    
  if(isCheckMode)
    return;
  
  if(istrue)
  {
    console.log('makeChange');
    var element = document.getElementById(id);
    element.classList.add("checked_img");
    checkedCount++;
    
    EnableCheckMode();
  }
}

document.addEventListener('touchmove', function(e) {
    istrue =false;
    clearTimeout(timer);
}, false);

function revert()
{
   istrue =false;
   clearTimeout(timer);
}

function goPost()
{
  var e = document.getElementById("selectDir");
  var destDir = e.value;
  if(destDir == '-filemove-')
  {
    return;
  }

  var list = document.getElementsByClassName('checked_img');
  if(list.length == 0)
  {
    return;
  }

  var name = [];
  for(var i = 0; i < list.length; i++)
  {
    name.push(list[i].name.substring(list[i].name.lastIndexOf('/')+1))
  }

  sourceDir = list[0].name.substring(list[0].name.indexOf('images'), list[0].name.lastIndexOf('/')+1)

  var form = document.createElement("form");
  form.setAttribute("method", "POST");
  form.setAttribute("action", "/api/input");
  form.setAttribute("target", "iframe1");

  //히든으로 값을 주입시킨다.
  var hiddenField = document.createElement("input");
  hiddenField.setAttribute("type", "hidden");
  hiddenField.setAttribute("name", "files");
  hiddenField.setAttribute("value", name);
  form.appendChild(hiddenField);

  hiddenField = document.createElement("input");
  hiddenField.setAttribute("type", "hidden");
  hiddenField.setAttribute("name", "dest");
  hiddenField.setAttribute("value", destDir);
  form.appendChild(hiddenField);

  hiddenField = document.createElement("input");
  hiddenField.setAttribute("type", "hidden");
  hiddenField.setAttribute("name", "source");
  hiddenField.setAttribute("value", sourceDir);
  form.appendChild(hiddenField);

  document.body.appendChild(form);
  form.submit();

  location.reload();
}

function removeSelect()
{
  const list = document.getElementsByClassName('checked_img');

  for(var i = 0; i != list.length;)
  {
    list[0].classList.remove("checked_img");
  }

  DisableCheckMode();
  checkedCount = 0;
}

function EnableCheckMode()
{
  isCheckMode = true;
  
  var div1 = document.getElementById("bot-check");
  var div2 = document.getElementById("bot-noncheck");
  div1.style.display = "block";
  div2.style.display = "none";
}

function DisableCheckMode()
{
  isCheckMode = false;
  
  var div1 = document.getElementById("bot-check");
  var div2 = document.getElementById("bot-noncheck");
  div1.style.display = "none";
  div2.style.display = "block";
}

dragElement(modal);

function dragElement(elmnt) {
  var pos1 = 0, pos3 = 0;
  modal.onmousedown = dragMouseDown;
	modal.onpointerdown = dragMouseDown;

  function dragMouseDown(e) {
//    e = e || window.event;
//    e.preventDefault();
    // get the mouse cursor position at startup:
    pos3 = e.clientX;
    pos4 = e.clientY;
    document.onmouseup = closeDragElement;
    document.ontouchend = closeDragElement;
    // call a function whenever the cursor moves:
    document.onmousemove = elementDrag;
    document.ontouchmove = elementDrag;
  }

  function elementDrag(e) {
    if(!imgMode)  
      return;

    pos1 = e.touches[0].pageX - pos3;
    pos3 = e.touches[0].pageX; 

    var style = window.getComputedStyle(modal);
    var matrix = new WebKitCSSMatrix(style.transform);

    elmnt.style.transform = "translate(" + (matrix.m41 + pos1) + "px)";
	
//	    if(e.type == 'touchstart' || e.type == 'touchmove' || e.type == 'touchend' || e.type == 'touchcancel'){
 //       var touch = e.originalEvent.touches[0] || e.originalEvent.changedTouches[0];
  //      x = touch.pageX;
   //     y = touch.pageY;
 //   } else if (e.type == 'mousedown' || e.type == 'mouseup' || e.type == 'mousemove' || e.type == 'mouseover'|| e.type=='mouseout' || e.type=='mouseenter' || e.type=='mouseleave') {
  //      x = e.clientX;
   //     y = e.clientY;
  //  }
  }

  function closeDragElement() {
    var style = window.getComputedStyle(modal);
    var matrix = new WebKitCSSMatrix(style.transform);
    if((matrix.m41) > (modal.offsetWidth/2))
    { 
      var nimg = document.getElementById('img' + (parseInt(img.name.substring(3))-1));
      if(nimg != null)
      {
        img.name = 'img' + (parseInt(img.name.substring(3))-1);
        img.src = nimg.name;
      }
    }
    else if((matrix.m41) < (modal.offsetWidth/2*-1))
    {
      var nimg = document.getElementById('img' + (parseInt(img.name.substring(3))+1));
      if(nimg != null)
      {
        img.name = 'img' + (parseInt(img.name.substring(3))+1);
        img.src = nimg.name;
      }
    }

    document.onmouseup = null;
    document.onmousemove = null;
    document.onpointerup = null;
    document.onpointermove = null;
    
    elmnt.style.transform = "translate(" + 0 + "px)";
  }
}

document.getElementById("prevPageBtn").addEventListener("click", function() {
  if(page <= 1)
    return;

  location.href=location.protocol + "//" + location.host + location.pathname + "?p=" + (page-1);
});

document.getElementById("nextPageBtn").addEventListener("click", function() {
if(lastPage)
  return;

  location.href=location.protocol + "//" + location.host + location.pathname + "?p=" + (page+1);
});