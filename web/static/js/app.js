
function postToServer(data) {
    $.ajax({url: "/actions",
        type: "POST",
        dataType: "json",
        contentType: "application/json; charset=utf-8",
        data: JSON.stringify(data),
        success: postCallback,
        error: function(xhr, status, error) {
            console.log("POST error: " + error);
            location.reload(true);
        }
    });
}

function postCallback(data, status, xhr) {
    console.log("status " + status);
    location.reload(true);
}

function buttonPressed(treePos) {
        postToServer({tree_pos: treePos});
}

$(document).ready(function() {
    $('#btn-reload').bind("click",function(){
            location.reload(true);
     });
    $('#btn-plate-solve').bind("click",function(){
        postToServer({action: "PLATE_SOLVER"});
     });
    $('#btn-snapshot').bind("click",function(){
        postToServer({action: "SNAPSHOT"});
     });

});