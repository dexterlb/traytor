import QtQuick 2.0
import QtQuick.Controls 1.4
import QtQuick.Window 2.2

Window {
    id: window
    width: 800
    height: 600


    ListModel {
        id: workers
        ListElement {
            address: "hoth"
            activeColour: "green"
        }
        ListElement {
            address: "do"
            activeColour: "red"
        }
    }

    Rectangle {
        id: page
        width: window.width
        height: window.height
        color: "lightgray"


        Rectangle{
            color: "lightgray"
            id : frame
            z: 15
            anchors.margins: 15
            anchors.top: page.top
            anchors.bottom: workerGrid.top
            anchors.left: page.left
            anchors.right: page.right
        }


        Image {
            id : image
            z: 20
            source: "foo.png"
            fillMode: Image.PreserveAspectFit
            anchors.left : frame.left
            anchors.top : frame.top
            width : frame.width
            height : frame.height
        }

        Rectangle {
            id: workerGrid
            color: "lightgray"
            anchors.bottom: newAddress.top
            anchors.left: page.left
            anchors.right: page.right
            anchors.margins: 15

            height: 100

            GridView {
                cellHeight: 30
                cellWidth: 150
                anchors.fill: parent
                model: workers
                delegate: Cell {
                    active: activeColour
                    workerAddress: address
                }
                boundsBehavior: Flickable.StopAtBounds
            }
        }

        Rectangle {
            id : newAddress
            anchors.bottom : page.bottom
            width : page.width - 70
            height : 30
            anchors.left : page.left
            function addWorker(){
                var workerAddress = inputAddress.text
                workers.append({
                    address: workerAddress,
                    activeColour: "green",
                })
                inputAddress.text = "";
            }


            TextInput {
                id : inputAddress
                anchors.fill : newAddress
                Keys.onReturnPressed: {
                newAddress.addWorker()
                boundsBehavior: Flickable.StopAtBounds
                }
            }

            Button {
                id : addWorker
                anchors.left : inputAddress.right
                width: page.width - newAddress.width
                height : newAddress.height
                text : "Add"

                MouseArea {
                    anchors.fill : parent
                    onClicked : newAddress.addWorker()
                }

            }
        }

    }
}
