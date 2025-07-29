package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateAuction_AutoClose(t *testing.T) {
	// Configurar um intervalo curto para o teste
	os.Setenv("AUCTION_INTERVAL", "100ms")
	defer os.Unsetenv("AUCTION_INTERVAL")

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:admin@localhost:27017/auctions?authSource=admin"))
	if err != nil {
		t.Skip("MongoDB não está disponível para teste")
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		t.Skip("MongoDB não está respondendo")
	}

	repo := NewAuctionRepository(client.Database("auctions"))

	// Cria o leilão
	auction, err := auction_entity.CreateAuction(
		"Test Product",
		"Electronics",
		"Test Description",
		auction_entity.New,
	)
	assert.Nil(t, err)

	err = repo.CreateAuction(ctx, auction)
	assert.Nil(t, err)

	// Aguarda o intervalo configurado + um pouco mais para garantir que o fechamento automático executou
	time.Sleep(200 * time.Millisecond)

	collection := client.Database("auctions").Collection("auctions")
	var result AuctionEntityMongo
	err = collection.FindOne(ctx, bson.M{"_id": auction.Id}).Decode(&result)

	if err != nil {
		t.Logf("Erro ao buscar leilão: %v", err)
		t.Skip("Não foi possível verificar o status do leilão")
	}

	// Verifica se o status mudou para Completed
	assert.Equal(t, auction_entity.Completed, result.Status, "O leilão deveria ter sido encerrado automaticamente")
}

func TestGetAuctionInterval(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected time.Duration
	}{
		{
			name:     "should return configured duration",
			envValue: "2m",
			expected: 2 * time.Minute,
		},
		{
			name:     "should return default when env not set",
			envValue: "",
			expected: 5 * time.Minute,
		},
		{
			name:     "should return default when env invalid",
			envValue: "invalid",
			expected: 5 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("AUCTION_INTERVAL", tt.envValue)
				defer os.Unsetenv("AUCTION_INTERVAL")
			} else {
				os.Unsetenv("AUCTION_INTERVAL")
			}

			result := getAuctionInterval()

			assert.Equal(t, tt.expected, result)
		})
	}
}
