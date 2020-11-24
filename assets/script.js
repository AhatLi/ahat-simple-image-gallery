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