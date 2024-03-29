var modal = document.getElementById("myModal");
var btn = document.getElementsByName("myBtn");
var img = document.getElementById("myImg");

var imgLayer = document.getElementById("modalImage");
var selectLayer = document.getElementById("modalSelect");
var configLayer = document.getElementById("modalConfig");
var searchLayer = document.getElementById("modalSearch");
var removeLayer = document.getElementById("modalRemove");

var timer;
var istrue = false;
var delay = 1000;
var isCheckMode = false;
var checkedCount = 0;
var imgMode = false;

function init() {  
  var pageSelect = document.getElementById("pageSelect");
  pageSelect.addEventListener('change', (event) => {
    var selectValue = pageSelect.options[pageSelect.selectedIndex].value;
    location.href=location.protocol + "//" + location.host + location.pathname + "?p=" + selectValue;
  });

  document.addEventListener('touchmove', function(e) {
      istrue =false;
      clearTimeout(timer);
  }, false);

  dragElement(modal);

  document.getElementsByName("imgCount")[0].value = contentCount;
  document.getElementsByName("imgSort")[0].value = contentSort;
}
init();

function setPageParam(search, p)
{
  var params = "?";
  var param = search.substr(1).split('&');
  for(var i = 0; i < param.length; i++)
  {
    var tmp = param[i].split('=');
    if(tmp[0] == 'p')
    {
      params += 'p=';
      params += p;
    }
    else
    {
      params += param[i]
    }
    if(i+1 != param.length)
    params += '&';
  }
}

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
    modalNone();
    imgLayer.style.display = "block";
    
    img.src = element.title;
    img.title = id;
    imgMode = true;
  }
}

function openModalSetting()
{
  modalNone();
  configLayer.style.display = "block";
}

function openModalSearch()
{
  modalNone();
  searchLayer.style.display = "block";
}

function openModalFileMove()
{
  modalNone();
  selectLayer.style.display = "block";
}

function openModalFileRemove()
{
  modalNone();
  removeLayer.style.display = "block";
}

function modalNone()
{
  modal.style.display = "block";
  imgLayer.style.display = "none";
  selectLayer.style.display = "none";
  configLayer.style.display = "none";
  searchLayer.style.display = "none";
  removeLayer.style.display = "none";
}

window.onclick = function(event) {
  if (event.target == modal) {
    modal.style.display = "none";
    imgMode = false;
  }
}

function fileRemove()
{
  var list = document.getElementsByClassName('checked_img');
  if(list.length == 0)
  {
    return;
  }

  var name = [];
  for(var i = 0; i < list.length; i++)
  {
    name.push(list[i].title.substring(list[i].title.lastIndexOf('/')+1))
  }

  removeDir = list[0].title.substring(list[0].title.indexOf('images'), list[0].title.lastIndexOf('/')+1)

  var form = document.createElement("form");
  form.setAttribute("method", "POST");
  form.setAttribute("action", "/api/delete");
  form.setAttribute("target", "tmpiframe");

  var hiddenField = document.createElement("input");
  hiddenField.setAttribute("type", "hidden");
  hiddenField.setAttribute("name", "files");
  hiddenField.setAttribute("value", name);
  form.appendChild(hiddenField);

  hiddenField = document.createElement("input");
  hiddenField.setAttribute("type", "hidden");
  hiddenField.setAttribute("name", "path");
  hiddenField.setAttribute("value", removeDir);
  form.appendChild(hiddenField);

  document.body.appendChild(form);
  form.submit();

  location.reload();
}

function longTouch(id)
{
   istrue = true;
   timer = setTimeout(function(){ imgCheck(id);},delay);
}

function imgCheck(id)
{
  if(timer)
    clearTimeout(timer);
    
  if(isCheckMode)
    return;
  
  if(istrue)
  {
    var element = document.getElementById(id);
    element.classList.add("checked_img");
    checkedCount++;
    
    EnableCheckMode();
  }
}

function revert()
{
   istrue =false;
   clearTimeout(timer);
}

function postFileMove()
{
  var destDir = document.getElementById("selectDir").value;
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
    name.push(list[i].title.substring(list[i].title.lastIndexOf('/')+1))
  }

  sourceDir = list[0].title.substring(list[0].title.indexOf('images'), list[0].title.lastIndexOf('/')+1)

  var form = document.createElement("form");
  form.setAttribute("method", "POST");
  form.setAttribute("action", "/api/move");
  form.setAttribute("target", "tmpiframe");

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

function postConfig()
{
  configForm = document.getElementById("configForm");
  configForm.submit();

  location.reload();
}

function getFileSearch()
{
  searchForm = document.getElementById("searchForm");
  searchForm.submit();
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
      var nimg = document.getElementById('img' + (parseInt(img.title.substring(3))-1));
      if(nimg != null)
      {
        img.title = 'img' + (parseInt(img.title.substring(3))-1);
        img.src = nimg.title;
      }
    }
    else if((matrix.m41) < (modal.offsetWidth/2*-1))
    {
      var nimg = document.getElementById('img' + (parseInt(img.title.substring(3))+1));
      if(nimg != null)
      {
        img.title = 'img' + (parseInt(img.title.substring(3))+1);
        img.src = nimg.title;
      }
    }

    document.onmouseup = null;
    document.onmousemove = null;
    document.onpointerup = null;
    document.onpointermove = null;
    
    elmnt.style.transform = "translate(" + 0 + "px)";
  }
  
	const index = 1;

  window.showGallery = function(index, id){
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
      const options = {
        class: "only-this-gallery",
        index: index+1,
        animation: ["slide", "fade"],
        autohide: "all",
        control: ["page", "theme", "autofit", "fullscreen", "zoom-in", "zoom-out", "close", "download",  "prev", "next"],
        fit: "contain"
      };
      Spotlight.show(gallery, options);
    }
  }
}