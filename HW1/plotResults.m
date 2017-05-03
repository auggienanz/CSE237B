% CSE 237B, Lab 1, Part 1
% August Nanz

% Import data
data = csvread('./client/result.csv');
time = (data(:,1) - data(1,1))./1e9;
ntp = data(:,2)./1e6;
offset = data(:,3)./1e6;
lambda = data(:,4)./1e6;

% Plot RTT, offset, and bound on error
figure;
plot(time, offset);
hold on
plot(time, offset - lambda,':');
plot(time, offset + lambda,':');
plot(time, ntp);
legend('Offset','Lower Bound','Upper Bound', 'NTP');
ylabel('Offset (ms)');
xlabel('Time (s)');
