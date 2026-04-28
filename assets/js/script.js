var modal = document.getElementById("myModal");
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
  if (pageSelect) {
    pageSelect.addEventListener("change", function () {
      var selectValue = pageSelect.options[pageSelect.selectedIndex].value;
      var url = new URL(window.location.href);
      url.searchParams.set("p", selectValue);
      window.location.href = url.toString();
    });
  }

  document.addEventListener("touchmove", function () {
    istrue = false;
    clearTimeout(timer);
  }, false);

  dragElement(modal);

  if (window.contentConfig) {
    document.getElementsByName("imgCount")[0].value = String(contentConfig.count);
    document.getElementsByName("imgSort")[0].value = contentConfig.sort;
    document.getElementsByName("mobileColumns")[0].value = String(contentConfig.mobileColumns);
  }
}
init();

function getSelectedImages() {
  return Array.from(document.getElementsByClassName("checked_img"));
}

function thumbClick(id) {
  var element = document.getElementById(id);

  if (isCheckMode) {
    toggleChecked(element);
    return;
  }

  modalNone();
  imgLayer.style.display = "block";
  img.src = element.dataset.fileUrl;
  img.title = id;
  imgMode = true;
}

function toggleChecked(element) {
  if (element.classList.contains("checked_img")) {
    element.classList.remove("checked_img");
    checkedCount--;
  } else {
    element.classList.add("checked_img");
    checkedCount++;
  }

  if (checkedCount <= 0) {
    checkedCount = 0;
    DisableCheckMode();
  }
}

function openModalSetting() {
  modalNone();
  configLayer.style.display = "block";
}

function openModalSearch() {
  modalNone();
  searchLayer.style.display = "block";
}

function openModalFileMove() {
  modalNone();
  selectLayer.style.display = "block";
}

function openModalFileRemove() {
  modalNone();
  removeLayer.style.display = "block";
}

function modalNone() {
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

async function postAction(url, payload) {
  var response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded;charset=UTF-8"
    },
    body: new URLSearchParams(payload)
  });

  if (!response.ok) {
    var message = await response.text();
    throw new Error(message || "request failed");
  }

  return response.json();
}

async function fileRemove() {
  var list = getSelectedImages();
  if (list.length === 0) {
    return;
  }

  var names = list.map(function(item) { return item.dataset.fileName; });
  var removeDir = list[0].dataset.dirPath;

  try {
    await postAction("/api/delete", {
      files: names.join(","),
      path: removeDir
    });
    window.location.reload();
  } catch (error) {
    alert(error.message);
  }
}

function longTouch(id) {
  istrue = true;
  timer = setTimeout(function() { imgCheck(id); }, delay);
}

function imgCheck(id) {
  if (timer) {
    clearTimeout(timer);
  }

  if (isCheckMode) {
    return;
  }

  if (istrue) {
    var element = document.getElementById(id);
    element.classList.add("checked_img");
    checkedCount++;
    EnableCheckMode();
  }
}

function revert() {
  istrue = false;
  clearTimeout(timer);
}

async function postFileMove() {
  var destDir = document.getElementById("selectDir").value;
  var list = getSelectedImages();
  if (list.length === 0) {
    return;
  }

  var names = list.map(function(item) { return item.dataset.fileName; });
  var sourceDir = list[0].dataset.dirPath;

  try {
    await postAction("/api/move", {
      files: names.join(","),
      dest: destDir,
      source: sourceDir
    });
    window.location.reload();
  } catch (error) {
    alert(error.message);
  }
}

async function postConfig() {
  var configForm = document.getElementById("configForm");
  var formData = new FormData(configForm);

  try {
    await postAction("/api/config", {
      imgCount: formData.get("imgCount"),
      imgSort: formData.get("imgSort"),
      mobileColumns: formData.get("mobileColumns")
    });
    window.location.reload();
  } catch (error) {
    alert(error.message);
  }
}

function removeSelect() {
  var list = getSelectedImages();
  list.forEach(function(item) {
    item.classList.remove("checked_img");
  });

  DisableCheckMode();
  checkedCount = 0;
}

function EnableCheckMode() {
  isCheckMode = true;

  var div1 = document.getElementById("bot-check");
  var div2 = document.getElementById("bot-noncheck");
  div1.style.display = "block";
  div2.style.display = "none";
}

function DisableCheckMode() {
  isCheckMode = false;

  var div1 = document.getElementById("bot-check");
  var div2 = document.getElementById("bot-noncheck");
  div1.style.display = "none";
  div2.style.display = "block";
}

function dragElement(elmnt) {
  var pos1 = 0;
  var pos3 = 0;
  modal.onmousedown = dragMouseDown;
  modal.onpointerdown = dragMouseDown;

  function dragMouseDown(e) {
    pos3 = e.clientX;
    document.onmouseup = closeDragElement;
    document.ontouchend = closeDragElement;
    document.onmousemove = elementDrag;
    document.ontouchmove = elementDrag;
  }

  function elementDrag(e) {
    if (!imgMode || !e.touches || e.touches.length === 0) {
      return;
    }

    pos1 = e.touches[0].pageX - pos3;
    pos3 = e.touches[0].pageX;

    var style = window.getComputedStyle(modal);
    var matrix = new WebKitCSSMatrix(style.transform);

    elmnt.style.transform = "translate(" + (matrix.m41 + pos1) + "px)";
  }

  function closeDragElement() {
    var style = window.getComputedStyle(modal);
    var matrix = new WebKitCSSMatrix(style.transform);
    if ((matrix.m41) > (modal.offsetWidth / 2)) {
      var nimg = document.getElementById("img" + (parseInt(img.title.substring(3), 10) - 1));
      if (nimg != null) {
        img.title = "img" + (parseInt(img.title.substring(3), 10) - 1);
        img.src = nimg.dataset.fileUrl;
      }
    } else if ((matrix.m41) < (modal.offsetWidth / 2 * -1)) {
      var nextImg = document.getElementById("img" + (parseInt(img.title.substring(3), 10) + 1));
      if (nextImg != null) {
        img.title = "img" + (parseInt(img.title.substring(3), 10) + 1);
        img.src = nextImg.dataset.fileUrl;
      }
    }

    document.onmouseup = null;
    document.onmousemove = null;
    document.onpointerup = null;
    document.onpointermove = null;

    elmnt.style.transform = "translate(0px)";
  }

  window.showGallery = function(index, element) {
    if (isCheckMode) {
      toggleChecked(element);
      return;
    }

    var options = {
      class: "only-this-gallery",
      index: index + 1,
      animation: ["slide", "fade"],
      autohide: "all",
      control: ["page", "theme", "autofit", "fullscreen", "zoom-in", "zoom-out", "close", "download", "prev", "next"],
      fit: "contain"
    };
    Spotlight.show(gallery, options);
  }
}
