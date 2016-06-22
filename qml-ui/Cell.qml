import QtQuick 2.6


Item {
    id: container
    property alias active: circle.color
    property alias workerAddress: worker.text
    signal clicked(color active)
    function changeColor(active) {
        if (active == "#ff0000")
            return "green";
        else
            return "red";
    }

    width: 100; height: 25
    Text {
        id: worker
        text : "gs, gs"
        anchors.left : circle.right
        y : 5
        height : 25
        width : 75
        anchors.margins : 10
    }

    Rectangle {
     id : circle
     width: parent.width<parent.height?parent.width:parent.height
     height: width
     border.color: "white"
     border.width: 1
     anchors.left : parent.left
     radius: width*0.5
}

    MouseArea {
        anchors.fill: parent
        onClicked: container.active = changeColor(container.active)
    }
}