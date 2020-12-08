var modal = document.getElementById("myModal");
var btn = document.getElementsByName("myBtn");
var img = document.getElementById("myImg");
var imgLayer = document.getElementById("modalImage");
var selectLayer = document.getElementById("modalSelect");

var timer;
var istrue = false;
var delay = 1000;
var isCheckMode = false;
var checkedCount = 0;

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
    img.src = element.name;
  }

  console.log('count : ' + checkedCount);
}

function movdImage()
{
  modal.style.display = "block";
  imgLayer.style.display = "none";
  selectLayer.style.display = "block";
}

window.onclick = function(event) {
  if (event.target == modal) {
    modal.style.display = "none";
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
  //  e.preventDefault();
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
  const list = document.getElementsByClassName('checkedItem');

  for(var i = 0; i < list.length; i++)
  {
    list[i].style.display = "block";
  }
}

function DisableCheckMode()
{
  isCheckMode = false;
  const list = document.getElementsByClassName('checkedItem');

  for(var i = 0; i < list.length; i++)
  {
    list[i].style.display = "none";
  }
}

var element = document.getElementById("modalImage");
dragElement(element);

function dragElement(elmnt) {
  var pos1 = 0, pos3 = 0;
  element.onmousedown = dragMouseDown;
	element.onpointerdown = dragMouseDown;

  function dragMouseDown(e) {
    e = e || window.event;
    e.preventDefault();
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

    pos1 = e.touches[0].pageX - pos3;
    pos3 = e.touches[0].pageX; 

    var style = window.getComputedStyle(element);
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

    document.onmouseup = null;
    document.onmousemove = null;
    document.onpointerup = null;
    document.onpointermove = null;
    
    elmnt.style.transform = "translate(" + 0 + "px)";

    var style = window.getComputedStyle(element);
    var matrix = new WebKitCSSMatrix(style.transform);
    if((matrix.m41 + pos1) > (element.offsetWidth/3*2))
    {
      //다음이미지
    }
    else if((matrix.m41 + pos1) < (element.offsetWidth/3*2*-1))
    {
      //이전이미지
    }
  }
}