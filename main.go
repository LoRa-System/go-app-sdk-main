package main

/* 종속성 컨트롤 */
import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	ttnsdk "github.com/TheThingsNetwork/go-app-sdk"
	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/go-utils/log/apex"
	"github.com/TheThingsNetwork/go-utils/random"
	"github.com/TheThingsNetwork/ttn/core/types"
)

const (
	// 애플리케이션 이름 (상수로 선언)
	sdkClientName = "my-amazing-app"
)

/* 메인함수 */
func main() {

	// 로깅 : 애플리케이션 프로그램의 표준 출력 기록
	log := apex.Stdout() // We use a cli logger at Stdout
	ttnlog.Set(log)      // Set the logger as default for TTN

	// 환경변수 -> 애플리케이션 ID, 액세스 키 가져옴
	appID := os.Getenv("TTN_APP_ID")
	appAccessKey := os.Getenv("TTN_APP_ACCESS_KEY")

	// 디버깅으로 출력하는거임
	log.Debug(appID + " " + appAccessKey)

	// 공용 커뮤니티 네트워크에 연결할 새 SDK 구성
	config := ttnsdk.NewCommunityConfig(sdkClientName)
	config.ClientVersion = "2.0.5" // The version of the application

	/* 새 SDK 클라이언트 생성 */
	client := config.NewClient(appID, appAccessKey)
	// 프로그램 종료 전, 클라이언트 정리 유무 확인
	defer client.Close()

	log.Debug("여기가 에러냐?")

	/* 새 장치 추가 -> SDK에서 장치관리자를 호출 */
	devices, err := client.ManageDevices()
	if err != nil {
	  log.WithError(err).Fatalf("%s: could not read CA certificate file", sdkClientName)
	}

	/* 새 장치 정보 입력 부분 */
	dev := new(ttnsdk.Device)
	dev.AppID = appID
	dev.DevID = "my-new-device"
	dev.Description = "A new device in my amazing app"
	dev.AppEUI = types.AppEUI{0x70, 0xB3, 0xD5, 0x7E, 0xF0, 0x00, 0x00, 0x24} // Use the real AppEUI here
	dev.DevEUI = types.DevEUI{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08} // Use the real DevEUI here
	random.FillBytes(dev.AppKey[:])    // Generate a random AppKey

	// Set 네트워크 서버에서 장치 사용
	err = devices.Set(dev)
	if err != nil {
		log.WithError(err).Fatalf("%s: could not create device", sdkClientName)
	}

	// GET 서버에서 장치를 되돌릴 수 있음
	dev, err = devices.Get("my-new-device")
	if err != nil {
		log.WithError(err).Fatalf("%s: could not get device", sdkClientName)
	}

	// 새로 생성된 장치가 네트워크에 연결하 준비가 되면
	// 활성화에 가입 -> MQTT클라이언트 시작
	pubsub, err := client.PubSub()
	if err != nil {
		log.WithError(err).Fatalf("%s: could not get application pub/sub", sdkClientName)
	}

	// 애플리케이션 내의 모든 장치에 대한 활성화를 받기를 원하면 AllDevices()
	allDevicesPubSub := pubsub.AllDevices()

	// List the first 10 devices
	/*deviceList, err := devices.List(10, 0)
	if err != nil {
		log.WithError(err).Fatal("my-amazing-app: could not get devices")
	}
	log.Info("my-amazing-app: found devices")
	for _, device := range deviceList {
		fmt.Printf("- %s", device.DevID)
	}*/

	// 장치 활성화에 가입 가능 -> 장치 연결 성공 시 콘솔 출력
	activations, err := allDevicesPubSub.SubscribeActivations()
	if err != nil {
		log.WithError(err).Fatalf("%s: could not subscribe to activations", sdkClientName)
	}
	go func() {
		for activation := range activations {
			log.WithFields(ttnlog.Fields{
				"appEUI":  activation.AppEUI.String(),
				"devEUI":  activation.DevEUI.String(),
				"devAddr": activation.DevAddr.String(),
			}).Info("my-amazing-app: received activation")
		}
	}()

	// 활성화 구독 취소
	err = allDevicesPubSub.UnsubscribeActivations()
	if err != nil {
		log.WithError(err).Fatalf("%s: could not unsubscribe from activations", sdkClientName)
	}

	// 생성한 장치 선택
	myNewDevicePubSub := pubsub.Device("my-new-device")

	// 업링크 메시지 구독 -> 도착하면 콘솔에 출력
	uplink, err := myNewDevicePubSub.SubscribeUplink()
	if err != nil {
		log.WithError(err).Fatalf("%s: could not subscribe to uplink messages", sdkClientName)
	}
	go func() {
		for message := range uplink {
			hexPayload := hex.EncodeToString(message.PayloadRaw)
			log.WithField("data", hexPayload).Infof("%s: received uplink", sdkClientName)
		}
	}()

	/* 구독 취소 */
	/*err = myNewDevicePubSub.UnsubscribeUplink()
	if err != nil {
		log.WithError(err).Fatalf("%s: could not unsubscribe from uplink", sdkClientName)
	}*/

	/* 다운링크 메시지 게시 -> 페이로드 "AABC"를 포트 10으로 전송 */
	err = myNewDevicePubSub.Publish(&types.DownlinkMessage{
		PayloadRaw: []byte{0xaa, 0xbc},
		FPort:      10,
	})
	if err != nil {
		log.WithError(err).Fatalf("%s: could not schedule downlink message", sdkClientName)
	}
}
