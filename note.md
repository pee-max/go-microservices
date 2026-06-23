# grpc

## #1 main

```golang
var grpcAddr = ":3301" //just a example

func main(){
    //create context
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    //listen for the shutdown signal
    go func(){
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt) //Ctrl+C (os.Interrupt)  kill (syscall.SIGTERM)
        <- sigCh //Block until a signal is received
        cancel()
    }()
    //listen to the port
    lis,err := net.Listen("tcp", grpcAddr)
    if err != nil {
        log.Fatalf("err...")
    }
    //start grpc server handler
    service := NewServer()
    grpcService := grpcServer.NewServer()
    NewGRPCHandler(grpcService, service)
    //log server startup
    log.Printf("starting grpc server server on port: %v", lis.Addr().String())
    //start grpc business server
    go func(){
        if err := grpcService.Serve(lis); err != nil {
            log.Printf("err...")
            cancel()
        }
    }()
    //block until shutdown signal
    <- ctx.Done()
    log.Println("shutting down the server...")
	grpcService.GracefulStop() 
}
```

> 优雅关闭

## #2 grpcHandler

```golang
type grpcHandler struct {
	pb.UnimplementedGrpcServiceServer
	service *Service
}

func NewGRPCHandler(grpc *grpc.Server,service *Service){
    handler := &grpcHandler{
        service: service,
    }
    pb.Server(grpc,handler)
}

//the services...
func (h *grpcHandler)...{
    
}
```

## #3 Service

```golang
type Service struct {
    //anything that you need...
}

func NewService() *Service{
    return &Service{...}
}
```

## #4 grpcClient

```golang
//"grpc" can be repalced to the service name 
type grpcServiceClient struct{
    client *pb.GrpcServiceClient
    conn *grpc.ClientConn
}

func NewGrpcServiceClient() (*grpcServiceClient, error){
    grpcServiceURL := os.Getenv("GRPC_Service_URL")
    if grpcServiceURL == "" {
        grpcServiceURL = "..."
    }
    //create client
    conn, err := grpc.NewClient(grpcServiceURL,grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, err
    }
    client := pb.NewGrpcServiceClient(conn)
    return &grpcServiceClient {
        client: client,
        conn: conn
    }, nil
}

func (c *grpcServiceClient)Close(){
    if c.conn != nil {
        if err := c.conn.Close(); err != nil{
            return;
        }
    }
}
```

