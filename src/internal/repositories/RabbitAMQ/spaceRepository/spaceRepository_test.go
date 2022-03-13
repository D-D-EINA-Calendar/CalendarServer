package spaceRepository_test

import (
	"testing"

	"github.com/D-D-EINA-Calendar/CalendarServer/src/internal/core/domain"
	spaceRepo "github.com/D-D-EINA-Calendar/CalendarServer/src/internal/repositories/RabbitAMQ/spaceRepository"
	connection "github.com/D-D-EINA-Calendar/CalendarServer/src/pkg/connect"
	"github.com/D-D-EINA-Calendar/CalendarServer/src/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestReserve(t *testing.T) {
	//t.Skip() //remove for activating it
	assert := assert.New(t)
	rabbitConn, err := connection.New(constants.AMQPURL)
	assert.Equal(err, nil, "Shouldn't be an error")
	chMonitoring, err := rabbitConn.NewChannel()
	assert.Equal(err, nil, "Shouldn't be an error")
	spaceRepo := spaceRepo.New(chMonitoring)
	done, err := spaceRepo.Reserve(domain.Space{},domain.Hour{Hour: 12, Min: 30},domain.Hour{Hour: 13, Min: 30})
	assert.Equal(err, nil, "Shouldn't be an error")
	assert.Equal(done, true, "Should be true")
}

func TestReservBatch(t *testing.T) {
	//t.Skip() //remove for activating it
	assert := assert.New(t)
	rabbitConn, err := connection.New(constants.AMQPURL)
	assert.Equal(err, nil, "Shouldn't be an error")
	chMonitoring, err := rabbitConn.NewChannel()
	assert.Equal(err, nil, "Shouldn't be an error")
	spaceRepo := spaceRepo.New(chMonitoring)
	done, err := spaceRepo.ReserveBatch([]domain.Space{},domain.Hour{Hour: 12, Min: 30},domain.Hour{Hour: 13, Min: 30})
	assert.Equal(err, nil, "Shouldn't be an error")
	assert.Equal(done, true, "Should be true")
}