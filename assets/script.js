var modal = document.getElementById("myModal");
var btn = document.getElementsByName("myBtn");
var span = document.getElementsByClassName("close")[0];
var img = document.getElementById("myImg");

function thumbClick(path)
{
  modal.style.display = "block";
  img.src = path;
}

span.onclick = function() {
  modal.style.display = "none";
}

window.onclick = function(event) {
  if (event.target == modal) {
    modal.style.display = "none";
  }
}