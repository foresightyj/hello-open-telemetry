namespace dotnet_core_gateway
{
    public class ApiResponse<T>
    {
        public bool Success { get; set; }
        public string Error { get; set; }
        public T Data { get; set; }
    }
}