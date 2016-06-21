import QtQuick 2.0


Rectangle {
    id: page
    width: 320; height: 480
    color: "lightgray"

    Rectangle{
        id : textArea
        color : "white"
        anchors.horizontalCenter: page.horizontalCenter
        width : 200
        height : 30
        y : 30
    }

    TextInput{
        id : address
        anchors.fill : textArea
        width : textArea.width
        height : textArea.height
        
    }

    Grid {
        id: colorPicker
        x: 4; anchors.bottom: page.bottom; anchors.bottomMargin: 4
        rows: 2; columns: 3; spacing: 3

        Cell { active: "red"; workerAddress : ":1234"}
        Cell { active: "green"; workerAddress : "hoth:1234"}
    }

}

