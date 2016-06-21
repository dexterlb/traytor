import QtQuick 2.0

Item {
    id: container
    property alias active: rectangle.color
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
        anchors.left : rectangle.right
        height : 25
        width : 75
        anchors.margins : 10
    }

    Rectangle {
        id: rectangle
        border.color: "white"
        width : 25
        height : 25
        anchors.left : parent.left
    }

    MouseArea {
        anchors.fill: parent
        onClicked: container.active = changeColor(container.active)
    }
}