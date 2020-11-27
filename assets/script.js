var modal = document.getElementById("myModal");
var btn = document.getElementsByName("myBtn");
var span = document.getElementsByClassName("close")[0];
var img = document.getElementById("myImg");

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
      isCheckMode = false;
  }
  else
  {
    modal.style.display = "block";
    img.src = element.name;
  }
}

span.onclick = function() {
  modal.style.display = "none";
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
  
  if(istrue)
  {
    var element = document.getElementById(id);
    element.classList.add("checked_img");
    checkedCount++;
    
    console.log(element.name);
    isCheckMode = true;
  }
}
function revert()
{
   istrue =false;
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
}