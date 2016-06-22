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
            id: addresses
            x: 4; 
            property int items : 0
            anchors.top : frame.bottom
            anchors.margins : 3
            rows: 2
            columns: 4
            spacing: 3
        }

        Rectangle {
            id : newAddress
            anchors.bottom : page.bottom
            width : page.width - 70
            height : 30
            anchors.left : page.left
            function addWorker(){
                            var worker = inputAddress.text
                            var newCell = 'import QtQuick 2.6;
                                    Cell {
                                        active: "green"; workerAddress : "' +  String(worker) + '"
                                    }'
                            if (worker != "" && addresses.items <= 12){
                                if(addresses.items == addresses.rows * addresses.columns
                                    && 
                                addresses.items <= 8) {
                                    addresses.rows += 1;
                                    window.height += Math.max(0, addresses.rows - 2) * 100;
                                    newAddress.anchors.bottom = page.bottom;
                                }
                                var newObject = Qt.createQmlObject(
                                    newCell,
                                    addresses, "dynamicSnippet1"
                                );
                                addresses.items = addresses.items + 1;
                                inputAddress.text = "";
                            }
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
