import QtQuick 2.0
import QtQuick.Controls 1.4
import QtQuick.Window 2.2

Window {
    id: window
    width: image.sourceSize.width + 50
    height: image.sourceSize.height + 120
    Rectangle {
        id: page
        width: window.width
        height: window.height
        color: "lightgray"

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

        Component {
            id: workerListDelegate
            Cell {
                active: activeColour
                workerAddress: address
            }
        }

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

        Rectangle {
            color: "transparent"
            anchors.top: frame.bottom
            anchors.bottom: newAddress.top
            anchors.margins: 30
            width: page.width

            GridView {
                anchors.fill: parent
                model: workers
                delegate: workerListDelegate
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
            }


            TextInput {
                id : inputAddress
                anchors.fill : newAddress
                Keys.onReturnPressed: {
                newAddress.addWorker()
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
