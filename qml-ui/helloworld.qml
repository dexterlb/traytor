import QtQuick 2.0


Rectangle {
    id: page
    width: image.sourceSize.width + 50
    height: image.sourceSize.height + 100
    color: "lightgray"

    Rectangle{
        id : frame
        anchors.horizontalCenter: page.horizontalCenter
        width : image.sourceSize.width
        height : image.sourceSize.height
        y : 30
    }

    Image {
        id : image
        source: "foo.png"
        fillMode: Image.PreserveAspectFit
        anchors.left : frame.left
        anchors.top : frame.top
        width : frame.width
        height : frame.height
    }

    Grid {
        id: colorPicker
        x: 4; 
        anchors.top: frame.bottom
        anchors.margins : 3
        rows: 2
        columns: 1
        spacing: 3

        Cell { active: "red"; workerAddress : ":1234"}
        Cell { active: "green"; workerAddress : "hoth:1234"}
    }

}

