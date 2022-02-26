package horariosrv_test

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/D-D-EINA-Calendar/CalendarServer/src/internal/core/domain"
	"github.com/D-D-EINA-Calendar/CalendarServer/src/internal/core/services/horariosrv"
	mock_ports "github.com/D-D-EINA-Calendar/CalendarServer/src/mocks/mockups"
	"github.com/D-D-EINA-Calendar/CalendarServer/src/pkg/apperrors"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/streadway/amqp"
)

type mocks struct {
	horarioRepository *mock_ports.MockHorarioRepositorio
}

//Checks all the cases for the function GetAvailableHours of the service [horariosrv]
func TestGetAvailableHours(t *testing.T) {
	// · Mocks · //
	AvailableHours := simpleAvailableHours()
	ternaAsked := domain.Terna{
		Degree: "Ing.Informática",
		Year:   2,
		Group:  "1",
	}

	// · Test · //
	type args struct {
		terna domain.Terna
	}
	type want struct {
		result []domain.AvailableHours
		err    error
	}
	tests := []struct {
		name  string
		args  args
		want  want
		mocks func(m mocks)
	}{{
		name: "Should return avaiable hours correctly",
		args: args{terna: ternaAsked},
		want: want{result: AvailableHours},
		mocks: func(m mocks) {
			m.horarioRepository.EXPECT().GetAvailableHours(ternaAsked).Return(AvailableHours, nil)
		},
	},
		{
			name: "Should return error when not found",
			args: args{terna: ternaAsked},
			want: want{result: []domain.AvailableHours{}, err: apperrors.ErrNotFound},
			mocks: func(m mocks) {
				m.horarioRepository.EXPECT().GetAvailableHours(ternaAsked).Return([]domain.AvailableHours{}, apperrors.ErrNotFound)
			},
		},
		{
			name: "Should return error when [titulación] is empty",
			args: args{terna: domain.Terna{Year: 1, Group: "1"}},
			want: want{result: []domain.AvailableHours{}, err: apperrors.ErrInvalidInput},
			mocks: func(m mocks) {
				m.horarioRepository.EXPECT().GetAvailableHours(domain.Terna{Year: 1, Group: "1"}).Return([]domain.AvailableHours{}, apperrors.ErrInvalidInput)
			},
		},
		{
			name: "Should return error when [curso] is empty",
			args: args{terna: domain.Terna{Degree: "A", Group: "1"}},
			want: want{result: []domain.AvailableHours{}, err: apperrors.ErrInvalidInput},
			mocks: func(m mocks) {
				m.horarioRepository.EXPECT().GetAvailableHours(domain.Terna{Degree: "A", Group: "1"}).Return([]domain.AvailableHours{}, apperrors.ErrInvalidInput)
			},
		},
		{
			name: "Should return error when [Group] is empty",
			args: args{terna: domain.Terna{Degree: "A", Year: 1}},
			want: want{result: []domain.AvailableHours{}, err: apperrors.ErrInvalidInput},
			mocks: func(m mocks) {
				m.horarioRepository.EXPECT().GetAvailableHours(domain.Terna{Degree: "A", Year: 1}).Return([]domain.AvailableHours{}, apperrors.ErrInvalidInput)
			},
		},
	}
	// · Runner · //
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Prepare

			m := mocks{
				horarioRepository: mock_ports.NewMockHorarioRepositorio(gomock.NewController(t)),
			}

			tt.mocks(m)
			service := horariosrv.New(m.horarioRepository)

			//Execute
			result, err := service.GetAvailableHours(tt.args.terna)

			//Verify
			if tt.want.err != nil && err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}

			assert.Equal(t, tt.want.result, result)

		})

	}
}

//Returns a slice for having different cases in the tests
func simpleAvailableHours() []domain.AvailableHours {

	return []domain.AvailableHours{
		{

			Subject:        domain.Subject{Kind: domain.THEORICAL, Name: "IC"},
			RemainingHours: 5,
			MaxHours:       5,
			RemainingMin:   0,
			MaxMin:         0,
		},
		{
			Subject:        domain.Subject{Name: "Prog 1", Kind: domain.PRACTICES},
			RemainingHours: 2,
			MaxHours:       3,
			RemainingMin:   0,
			MaxMin:         0,
		},
	}

}

/////////////////////////////////////
// TEST UPDATE SCHEDULER ENTRIES ///
///////////////////////////////////

func TestUpdateEntries(t *testing.T) {
	// · Mocks · //

	// · Test · //
	type args struct {
		entries []domain.Entry
		terna   domain.Terna
	}
	type want struct {
		result string
		err    error
	}
	tests := []struct {
		name  string
		args  args
		want  want
		mocks func(m mocks)
	}{{
		name: "Create entry succed",
		args: args{entries: simpleEntries(), terna: simpleTerna()},
		want: want{result: currentDate(), err: nil},
		mocks: func(m mocks) {
			m.horarioRepository.EXPECT().CreateNewEntry(simpleEntries()[0]).Return(nil)
			m.horarioRepository.EXPECT().CreateNewEntry(simpleEntries()[1]).Return(nil)
			m.horarioRepository.EXPECT().DeleteAllEntries(simpleTerna()).Return(nil)
		},
	},
		{
			name: "Should return error if repository fails",
			args: args{entries: simpleEntries(), terna: simpleTerna()},
			want: want{result: "", err: apperrors.ErrSql},
			mocks: func(m mocks) {
				m.horarioRepository.EXPECT().CreateNewEntry(simpleEntries()[0]).Return(apperrors.ErrInternal)
				m.horarioRepository.EXPECT().DeleteAllEntries(simpleTerna()).Return(nil)
			},
		},
	}
	// · Runner · //
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			//Prepare
			m := mocks{
				horarioRepository: mock_ports.NewMockHorarioRepositorio(gomock.NewController(t)),
			}

			tt.mocks(m)
			service := horariosrv.New(m.horarioRepository)

			//Execute
			result, err := service.UpdateScheduler(tt.args.entries, tt.args.terna)

			//Verify operation succeded
			if tt.want.err != nil && err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}

			assert.Equal(t, tt.want.result, result)

			//Verify state changed

			//TODO use the getEntry function for verifying the entrie was created

		})

	}
}

func simpleEntries() []domain.Entry {
	return []domain.Entry{
		{
			Init: domain.NewHour(1, 1),
			End:  domain.NewHour(2, 2),
			Subject: domain.Subject{
				Kind: domain.THEORICAL,
				Name: "Prog 1",
			},
			Room:    domain.Room{Name: "1"},
			Weekday: domain.MOONDAY,
		},
		{
			Init: domain.NewHour(5, 0),
			End:  domain.NewHour(9, 0),
			Subject: domain.Subject{
				Kind: domain.THEORICAL,
				Name: "Prog 2",
			},
			Room:    domain.Room{Name: "2"},
			Weekday: domain.THUERSDAY,
		},
	}
}

func currentDate() string {

	return time.Now().Format("02/01/2006")

}

func simpleTerna() domain.Terna {
	return domain.Terna{
		Group:  "1",
		Year:   1,
		Degree: "Ing Informática",
	}
}

/////////////////////////////
// TEST LIST DEGREES      ///
/////////////////////////////

func TestListSubject(t *testing.T) {
	// · Mocks · //

	// · Test · //

	type want struct {
		result []domain.DegreeDescription
		err    error
	}
	tests := []struct {
		name  string
		want  want
		mocks func(m mocks)
	}{
		{
			name: "Succeded",
			want: want{result: simpleListDegreeDescriptions(), err: nil},
			mocks: func(m mocks) {
				m.horarioRepository.EXPECT().ListAllDegrees().Return(simpleListDegreeDescriptions(), nil)
			},
		},

		{
			name: "Repo failure",
			want: want{result: nil, err: apperrors.ErrInternal},
			mocks: func(m mocks) {
				m.horarioRepository.EXPECT().ListAllDegrees().Return(nil, apperrors.ErrInternal)
			},
		},
	}

	// · Runner · //
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			//Prepare
			m := mocks{
				horarioRepository: mock_ports.NewMockHorarioRepositorio(gomock.NewController(t)),
			}

			tt.mocks(m)
			service := horariosrv.New(m.horarioRepository)

			//Execute
			result, err := service.ListAllDegrees()

			//Verify operation succeded
			if tt.want.err != nil && err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}

			assert.Equal(t, tt.want.result, result)

		})

	}
}

func simpleListDegreeDescriptions() []domain.DegreeDescription {
	return []domain.DegreeDescription{
		{
			Name: "A",
			Groups: []domain.YearDescription{
				{Name: 1, Groups: []string{"a", "b"}},
				{Name: 2, Groups: []string{"a", "b"}},
			},
		},
		{
			Name: "B",
			Groups: []domain.YearDescription{
				{Name: 1, Groups: []string{"a"}},
				{Name: 2, Groups: []string{"a", "b", "c"}},
			},
		},
	}
}

////////////////////////
// TEST GET  ENTRIES ///
///////////////////////

//Checks all the cases for the function GetAvailableHours of the service [horariosrv]
func TestGetEntries(t *testing.T) {
	// · Mocks · //
	entries := simpleEntries()
	ternaAsked := domain.Terna{
		Degree: "Ing.Informática",
		Year:   2,
		Group:  "1",
	}

	// · Test · //
	type args struct {
		terna domain.Terna
	}
	type want struct {
		result []domain.Entry
		err    error
	}
	tests := []struct {
		name  string
		args  args
		want  want
		mocks func(m mocks)
	}{
		{
			name: "Should return entries correctly",
			args: args{terna: ternaAsked},
			want: want{result: entries},
			mocks: func(m mocks) {
				m.horarioRepository.EXPECT().GetEntries(ternaAsked).Return(entries, nil)
			},
		},
		{
			name: "Should return error when not found",
			args: args{terna: ternaAsked},
			want: want{result: []domain.Entry{}, err: apperrors.ErrNotFound},
			mocks: func(m mocks) {
				m.horarioRepository.EXPECT().GetEntries(ternaAsked).Return([]domain.Entry{}, apperrors.ErrNotFound)
			},
		},
		{
			name:  "Should return error when [titulación] is empty",
			args:  args{terna: domain.Terna{Year: 1, Group: "1"}},
			want:  want{result: []domain.Entry{}, err: apperrors.ErrInvalidInput},
			mocks: func(m mocks) {},
		},
		{
			name:  "Should return error when [Group] is empty",
			args:  args{terna: domain.Terna{Degree: "A", Year: 1}},
			want:  want{result: []domain.Entry{}, err: apperrors.ErrInvalidInput},
			mocks: func(m mocks) {},
		},
		{
			name:  "Should return error when [Year] is empty",
			args:  args{terna: domain.Terna{Degree: "A", Group: "1"}},
			want:  want{result: []domain.Entry{}, err: apperrors.ErrInvalidInput},
			mocks: func(m mocks) {},
		},
	}

	// · Runner · //
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Prepare

			m := mocks{
				horarioRepository: mock_ports.NewMockHorarioRepositorio(gomock.NewController(t)),
			}

			tt.mocks(m)
			service := horariosrv.New(m.horarioRepository)

			//Execute
			result, err := service.GetEntries(tt.args.terna)

			//Verify
			if tt.want.err != nil && err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}

			assert.Equal(t, tt.want.result, result)

		})

	}

}

func TestGetICS(t *testing.T) {
	// · Mocks · //
	entries := simpleEntries()
	ternaAsked := domain.Terna{
		Degree: "Ing.Informática",
		Year:   2,
		Group:  "1",
	}

	// · Test · //
	type args struct {
		terna domain.Terna
	}
	type want struct {
		result string
		err    error
	}
	tests := []struct {
		name  string
		args  args
		want  want
		mocks func(m mocks)
	}{
		{
			name: "Should return ICS correctly",
			args: args{terna: ternaAsked},
			want: want{result: ""},
			mocks: func(m mocks) {
				m.horarioRepository.EXPECT().GetEntries(ternaAsked).Return(entries, nil)
			},
		},
		{
			name: "Should return error when not found",
			args: args{terna: ternaAsked},
			want: want{result: "", err: apperrors.ErrSql},
			mocks: func(m mocks) {
				m.horarioRepository.EXPECT().GetEntries(ternaAsked).Return([]domain.Entry{}, apperrors.ErrNotFound)
			},
		},
		{
			name:  "Should return error when [titulación] is empty",
			args:  args{terna: domain.Terna{Year: 1, Group: "1"}},
			want:  want{result: "", err: apperrors.ErrInvalidInput},
			mocks: func(m mocks) {},
		},
		{
			name:  "Should return error when [Group] is empty",
			args:  args{terna: domain.Terna{Degree: "A", Year: 1}},
			want:  want{result: "", err: apperrors.ErrInvalidInput},
			mocks: func(m mocks) {},
		},
		{
			name:  "Should return error when [Year] is empty",
			args:  args{terna: domain.Terna{Degree: "A", Group: "1"}},
			want:  want{result: "", err: apperrors.ErrInvalidInput},
			mocks: func(m mocks) {},
		},
	}

	// · Runner · //
	for i , tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Prepare

			m := mocks{
				horarioRepository: mock_ports.NewMockHorarioRepositorio(gomock.NewController(t)),
			}

			tt.mocks(m)
			service := horariosrv.New(m.horarioRepository)

			//Execute
			result, err := service.GetICS(tt.args.terna)

			//Verify
			if tt.want.err != nil && err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}

			if i != 0 {
				assert.Equal(t, tt.want.result, result)
			} else {
				assert.NotEqual(t, "", result)
			}

		})

	}

}

func TestUpdateByCSV(t *testing.T) {
	// · Mocks · //

	// · Test · //
	type args struct {
		terna domain.Terna
	}
	type want struct {
		result bool
		err    error
	}
	tests := []struct {
		name  string
		args  args
		want  want
		mocks func(m mocks)
	}{
		{
			name: "Should Update correctly",
			args: args{},
			want: want{result: true},
			mocks: func(m mocks) {
				m.horarioRepository.EXPECT().CreateNewDegree(558,"Graduado en Ingeniería en Diseño Industrial y Desarrollo de Producto").Return(true, nil)
				m.horarioRepository.EXPECT().CreateNewYear(1,558).Return(true, nil)
				m.horarioRepository.EXPECT().RawExec("INSERT INTO `grupodocente` (`id`, `numero`, `idcurso`) VALUES ('55811','1','5581'), ('55812','2','5581')").Return(nil)
				m.horarioRepository.EXPECT().RawExec("INSERT INTO `asignatura` (`id`, `codigo`, `nombre`, `idT`) VALUES ('25802','25802','Informática','558')").Return(nil)
				m.horarioRepository.EXPECT().RawExec("INSERT INTO `hora` (`id`, `disponibles`, `totales`, `tipo`, `idasignatura`, `idgrupo`, `grupo`, `semana`) VALUES (NULL,'3500','3500','1','25802','55811','',''), (NULL,'1000','1000','2','25802','55811','1',''), (NULL,'1000','1000','2','25802','55811','2',''), (NULL,'1500','1500','3','25802','55811','1','a'), (NULL,'1500','1500','3','25802','55811','2','a'), (NULL,'1500','1500','3','25802','55811','3','a'), (NULL,'3500','3500','1','25802','55812','',''), (NULL,'1000','1000','2','25802','55812','1',''), (NULL,'1000','1000','2','25802','55812','2',''), (NULL,'1500','1500','3','25802','55812','1','a'), (NULL,'1500','1500','3','25802','55812','2','a'), (NULL,'1500','1500','3','25802','55812','3','a')").Return(nil)
			},
		},
		{
			name: "Subject creation fails",
			args: args{},
			want: want{result: false, err: apperrors.ErrSql},
			mocks: func(m mocks) {
				m.horarioRepository.EXPECT().CreateNewDegree(558,"Graduado en Ingeniería en Diseño Industrial y Desarrollo de Producto").Return(true, nil)
				m.horarioRepository.EXPECT().CreateNewYear(1,558).Return(true, nil)
				m.horarioRepository.EXPECT().RawExec("INSERT INTO `grupodocente` (`id`, `numero`, `idcurso`) VALUES ('55811','1','5581'), ('55812','2','5581')").Return(nil)
				m.horarioRepository.EXPECT().RawExec("INSERT INTO `asignatura` (`id`, `codigo`, `nombre`, `idT`) VALUES ('25802','25802','Informática','558')").Return(apperrors.ErrSql)
			},
		},
		{
			name: "Hour creation fails",
			args: args{},
			want: want{result: false, err: apperrors.ErrSql},
			mocks: func(m mocks) {
				m.horarioRepository.EXPECT().CreateNewDegree(558,"Graduado en Ingeniería en Diseño Industrial y Desarrollo de Producto").Return(true, nil)
				m.horarioRepository.EXPECT().CreateNewYear(1,558).Return(true, nil)
				m.horarioRepository.EXPECT().RawExec("INSERT INTO `grupodocente` (`id`, `numero`, `idcurso`) VALUES ('55811','1','5581'), ('55812','2','5581')").Return(nil)
				m.horarioRepository.EXPECT().RawExec("INSERT INTO `asignatura` (`id`, `codigo`, `nombre`, `idT`) VALUES ('25802','25802','Informática','558')").Return(nil)
				m.horarioRepository.EXPECT().RawExec("INSERT INTO `hora` (`id`, `disponibles`, `totales`, `tipo`, `idasignatura`, `idgrupo`, `grupo`, `semana`) VALUES (NULL,'3500','3500','1','25802','55811','',''), (NULL,'1000','1000','2','25802','55811','1',''), (NULL,'1000','1000','2','25802','55811','2',''), (NULL,'1500','1500','3','25802','55811','1','a'), (NULL,'1500','1500','3','25802','55811','2','a'), (NULL,'1500','1500','3','25802','55811','3','a'), (NULL,'3500','3500','1','25802','55812','',''), (NULL,'1000','1000','2','25802','55812','1',''), (NULL,'1000','1000','2','25802','55812','2',''), (NULL,'1500','1500','3','25802','55812','1','a'), (NULL,'1500','1500','3','25802','55812','2','a'), (NULL,'1500','1500','3','25802','55812','3','a')").Return(apperrors.ErrSql)
			},
		},
		{
			name: "Group creation fails",
			args: args{},
			want: want{result: false, err: apperrors.ErrSql},
			mocks: func(m mocks) {
				m.horarioRepository.EXPECT().CreateNewDegree(558,"Graduado en Ingeniería en Diseño Industrial y Desarrollo de Producto").Return(true, nil)
				m.horarioRepository.EXPECT().CreateNewYear(1,558).Return(true, nil)
				m.horarioRepository.EXPECT().RawExec("INSERT INTO `grupodocente` (`id`, `numero`, `idcurso`) VALUES ('55811','1','5581'), ('55812','2','5581')").Return(apperrors.ErrSql)
			},
		},
	}

	// · Runner · //
	for i , tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Prepare

			m := mocks{
				horarioRepository: mock_ports.NewMockHorarioRepositorio(gomock.NewController(t)),
			}

			tt.mocks(m)
			service := horariosrv.New(m.horarioRepository)
			content, err := ioutil.ReadFile("../../../../pkg/csv/Listado207_1Asig.csv")
			contentString := string(content)
			//Execute
			result, err := service.UpdateByCSV(contentString)

			//Verify
			if tt.want.err != nil && err != nil {
				assert.Equal(t, tt.want.err.Error(), err.Error())
			}

			if i != 0 {
				assert.Equal(t, tt.want.result, result)
			} else {
				assert.NotEqual(t, "", result)
			}

		})

	}

}

func TestRMQ(t *testing.T) {
	url := "amqps://cnvzbkyj:zrT84snzxNyFwAZl1MV2vI9Gg8OtjiRV@whale.rmq.cloudamqp.com/cnvzbkyj"
    connection, _ := amqp.Dial(url)
    defer connection.Close()
    go func(con *amqp.Connection) {
        channel, _ := connection.Channel()
        defer channel.Close()
        durable, exclusive := false, false
        autoDelete, noWait := true, true
        q, _ := channel.QueueDeclare("test", durable, autoDelete, exclusive, noWait, nil)
        channel.QueueBind(q.Name, "#", "amq.topic", false, nil)
        autoAck, exclusive, noLocal, noWait := false, false, false, false
        messages, _ := channel.Consume(q.Name, "", autoAck, exclusive, noLocal, noWait, nil)
        multiAck := false
        for msg := range messages {
            fmt.Println("Body:", string(msg.Body), "Timestamp:", msg.Timestamp)
            msg.Ack(multiAck)
        }
    }(connection)

    go func(con *amqp.Connection) {
        timer := time.NewTicker(1 * time.Second)
        channel, _ := connection.Channel()

        for t := range timer.C {
            msg := amqp.Publishing{
                DeliveryMode: 1,
                Timestamp:    t,
                ContentType:  "text/plain",
                Body:         []byte("Hello world"),
            }
            mandatory, immediate := false, false
            channel.Publish("amq.topic", "ping", mandatory, immediate, msg)
        }
    }(connection)

    select {}
}